package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/sync/errgroup"

	openapispec "github.com/go-openapi/spec"

	"github.com/magodo/ghwalk"
	"github.com/magodo/terraform-provider-azurerm-insight/pkg/core"
)

type SWGResourceProviders map[string]*SWGResourceProvider

func (providers SWGResourceProviders) Copy() SWGResourceProviders {
	out := map[string]*SWGResourceProvider{}
	for k, v := range providers {
		newV := v.Copy()
		out[k] = &newV
	}
	return out
}

func (providers SWGResourceProviders) CopyRPFrom(oproviders SWGResourceProviders, rpName string) {
	if rp, ok := oproviders[rpName]; ok {
		newRP := rp.Copy()
		providers[rpName] = &newRP
	}
	return
}

func (providers SWGResourceProviders) CopyAPIFrom(oproviders SWGResourceProviders, rpName, apiVersion string) {
	rp, ok := oproviders[rpName]
	if !ok {
		return
	}
	api, ok := rp.Apis[apiVersion]
	if !ok {
		return
	}

	if _, ok := providers[rpName]; !ok {
		providers[rpName] = &SWGResourceProvider{
			SwaggerRelPath: rp.SwaggerRelPath,
			Apis:           SWGResourceProviderAPIs{},
		}
	}

	newApi := api.Copy()
	providers[rpName].Apis[apiVersion] = &newApi
	return
}

func (providers SWGResourceProviders) CopySchemaFrom(oproviders SWGResourceProviders, rpName, apiVersion, schemaName string) {
	rp, ok := oproviders[rpName]
	if !ok {
		return
	}
	api, ok := rp.Apis[apiVersion]
	if !ok {
		return
	}
	schema, ok := api.Schemas[schemaName]
	if !ok {
		return
	}

	if _, ok := providers[rpName]; !ok {
		providers[rpName] = &SWGResourceProvider{
			SwaggerRelPath: rp.SwaggerRelPath,
			Apis:           SWGResourceProviderAPIs{},
		}
	}

	if _, ok := providers[rpName].Apis[apiVersion]; !ok {
		providers[rpName].Apis[apiVersion] = &SWGResourceProviderAPI{
			SwaggerRelPath: api.SwaggerRelPath,
			Schemas:        SWGSchemas{},
		}
	}

	// Note that the SWGSchema is shared (i.e. not copied)
	providers[rpName].Apis[apiVersion].Schemas[schemaName] = schema
	return
}

type SWGResourceProvider struct {
	SwaggerRelPath string
	Apis           SWGResourceProviderAPIs
}

func (v SWGResourceProvider) Copy() SWGResourceProvider {
	return SWGResourceProvider{
		SwaggerRelPath: v.SwaggerRelPath,
		Apis:           v.Apis.Copy(),
	}
}

type SWGResourceProviderAPIs map[string]*SWGResourceProviderAPI

func (apis SWGResourceProviderAPIs) Copy() SWGResourceProviderAPIs {
	out := map[string]*SWGResourceProviderAPI{}
	for k, v := range apis {
		newV := v.Copy()
		out[k] = &newV
	}
	return out
}

type SWGResourceProviderAPI struct {
	SwaggerRelPath string
	Schemas        SWGSchemas
}

func (v SWGResourceProviderAPI) Copy() SWGResourceProviderAPI {
	return SWGResourceProviderAPI{
		SwaggerRelPath: v.SwaggerRelPath,
		Schemas:        v.Schemas.ShallowCopy(),
	}
}

type SWGSchemas map[string]*SWGSchema

func (v SWGSchemas) ShallowCopy() SWGSchemas {
	out := map[string]*SWGSchema{}
	for k, v := range v {
		out[k] = v
	}
	return out
}

type SWGSchema struct {
	core.SWGSchema
}

type SWGSchemaAddr struct {
	raw                    core.SWGSchemaAddr
	ResourceProvider       string
	MidPathSegment         string // e.g. resource-manager/Microsoft.Foobar/stable
	ApiVersion             string
	ResourceCollectionName string
	SchemaName             string
}

func ParseSWGSchemaAddr(addr core.SWGSchemaAddr) SWGSchemaAddr {
	swaggerRelPath, schemaName := addr.SwaggerRelPath(), addr.SchemaName()
	segments := strings.Split(swaggerRelPath, "/")
	rp, mid, api, resourceCollectionName := segments[0], segments[1:len(segments)-2], segments[len(segments)-2], strings.TrimSuffix(segments[len(segments)-1], ".json")
	return SWGSchemaAddr{
		raw:                    addr,
		ResourceProvider:       rp,
		MidPathSegment:         strings.Join(mid, "/"),
		ApiVersion:             api,
		ResourceCollectionName: resourceCollectionName,
		SchemaName:             schemaName,
	}
}

func (addr SWGSchemaAddr) RelPathToRP() string {
	return addr.ResourceProvider
}

func (addr SWGSchemaAddr) RelPathToApiVersion() string {
	return fmt.Sprintf("%s/%s/%s", addr.RelPathToRP(), addr.MidPathSegment, addr.ApiVersion)
}

func (addr SWGSchemaAddr) RelPathToSwaggerFile() string {
	return fmt.Sprintf("%s/%s.json", addr.RelPathToApiVersion(), addr.ResourceCollectionName)
}

// NewSWGResourceProviders convert the core.SWGSchemas, whose key is swagger file + schema name,
// into a hierarchy of structures mapping to the Azure concept, beginning from the resource provider level.
func NewSWGResourceProviders(swgschemas core.SWGSchemas) SWGResourceProviders {
	out := map[string]*SWGResourceProvider{}
	for addr, swgschema := range swgschemas.GetAll() {
		addr := ParseSWGSchemaAddr(addr)

		if _, ok := out[addr.ResourceProvider]; !ok {
			out[addr.ResourceProvider] = &SWGResourceProvider{SwaggerRelPath: addr.RelPathToRP(), Apis: map[string]*SWGResourceProviderAPI{}}
		}

		if _, ok := out[addr.ResourceProvider].Apis[addr.ApiVersion]; !ok {
			out[addr.ResourceProvider].Apis[addr.ApiVersion] = &SWGResourceProviderAPI{SwaggerRelPath: addr.RelPathToApiVersion(), Schemas: map[string]*SWGSchema{}}
		}

		out[addr.ResourceProvider].Apis[addr.ApiVersion].Schemas[addr.SchemaName] = &SWGSchema{*swgschema}
	}

	return out
}

// CompleteSWGResourceProvidersViaGithubAPI completes the swagger resource providers by querying swagger spec repo via Github.
// For each (RP,API Version), searching for all the swagger spec files to collect all the schemas that belongs to
// the "in-body" parameter of an endpoint which has PUT and DELETE methods.
func (swgrps SWGResourceProviders) CompleteSWGResourceProvidersViaGithubAPI(ctx context.Context, options *ghwalk.WalkOptions) error {
	log.Println("Start to complete SWGResourceProviders")
	const (
		swaggerRepoOwner        = "Azure"
		swaggerRepoRepo         = "azure-rest-api-specs"
		swaggerRepoSpecBasePath = "specification"
		swaggerRepoBaseURI      = "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/master/specification"
	)
	for rpName, rp := range swgrps {
		for apiName, api := range rp.Apis {

			schemaFolderPattern := regexp.MustCompile(fmt.Sprintf(`^%s(/resource-manager(/Microsoft.\w+(/(preview|stable)(/%s)?)?)?)?$`, rpName, apiName))
			schemaPattern := regexp.MustCompile(fmt.Sprintf(`^%s/resource-manager/Microsoft.\w+/(preview|stable)/%s/\w+.json$`, rpName, apiName))
			//if err := ghwalk.Walk(ctx, swaggerRepoOwner, swaggerRepoRepo, "specification/network/resource-manager/Microsoft.Network/stable/2020-05-01/azureFirewall.json", options,
			if err := ghwalk.Walk(ctx, swaggerRepoOwner, swaggerRepoRepo, path.Join(swaggerRepoSpecBasePath, rpName), options,
				func(p string, info *ghwalk.FileInfo, err error) error {
					relPath := strings.TrimPrefix(p, swaggerRepoSpecBasePath+"/")
					log.Printf("Searching Swaggers in %s...\n", relPath)
					if err != nil {
						return err
					}
					if info == nil {
						return nil
					}

					if info.IsDir() {
						return nil
					}

					schemas, err := collectAllTFCandidateSchemas(swaggerRepoBaseURI, relPath)
					if err != nil {
						return err
					}

					for _, schema := range schemas {
						if _, ok := api.Schemas[schema.Name]; !ok {
							api.Schemas[schema.Name] = &schema
						}
					}

					return nil
				},
				func(p string, info *ghwalk.FileInfo) bool {
					relPath := strings.TrimPrefix(p, swaggerRepoSpecBasePath+"/")

					// Skip directories not match the schema folder patterns
					if info.IsDir() {
						if !schemaFolderPattern.MatchString(relPath) {
							log.Printf("Skip directory %s!\n", relPath)
							return true
						}
						return false
					}

					// Skip files not match the schema file patterns
					if !schemaPattern.MatchString(relPath) {
						log.Printf("Skip file %s!\n", relPath)
						return true
					}
					return false
				}); err != nil {
				return err
			}
		}
	}

	return nil
}

// CompleteSWGResourceProvidersViaLocalFS is similar to the CompleteSWGResourceProvidersViaGithubAPI, except it walks the Azure Swagger
// repo on local FS.
func (swgrps SWGResourceProviders) CompleteSWGResourceProvidersViaLocalFS(swaggerRepoSpecBasePath string) error {
	g := new(errgroup.Group)
	for rpName, rp := range swgrps {
		for apiName, api := range rp.Apis {

			// Copy the variables which will be used in the goroutine's closure
			// Especiall,y the api.Schemas should be shallow copied so that the modifications will be reflected to the swgrps
			apiSchemas := api.Schemas
			rpName := rpName
			apiName := apiName

			g.Go(func() error {
				schemaFolderPattern := regexp.MustCompile(fmt.Sprintf(`^%[1]s(%[3]sresource-manager(%[3]sMicrosoft.\w+(%[3]s(preview|stable)(%[3]s%[2]s)?)?)?)?$`, rpName, apiName, regexp.QuoteMeta(string(os.PathSeparator))))
				schemaPattern := regexp.MustCompile(fmt.Sprintf(`^%[1]s%[3]sresource-manager%[3]sMicrosoft.\w+%[3]s(preview|stable)%[3]s%[2]s%[3]s\w+.json$`, rpName, apiName, regexp.QuoteMeta(string(os.PathSeparator))))
				if err := filepath.Walk(path.Join(swaggerRepoSpecBasePath, rpName),
					func(p string, info os.FileInfo, err error) error {
						p, _ = filepath.Abs(p)
						swaggerRepoSpecBasePath, _ = filepath.Abs(swaggerRepoSpecBasePath)
						relPath := strings.TrimPrefix(p, swaggerRepoSpecBasePath+string(os.PathSeparator))
						log.Printf("Searching Swaggers in %s...\n", relPath)
						if err != nil {
							return err
						}
						if info == nil {
							return nil
						}

						// Skip directories not match the schema folder patterns
						if info.IsDir() {
							if !schemaFolderPattern.MatchString(relPath) {
								log.Printf("Skip directory %s!\n", relPath)
								return filepath.SkipDir
							}
							return nil
						}

						// Skip files not match the schema file patterns
						if !schemaPattern.MatchString(relPath) {
							log.Printf("Skip file %s!\n", relPath)
							return nil
						}

						schemas, err := collectAllTFCandidateSchemas(swaggerRepoSpecBasePath, relPath)
						if err != nil {
							return err
						}

						for _, schema := range schemas {
							if _, ok := apiSchemas[schema.Name]; !ok {
								apiSchemas[schema.Name] = &schema
							}
						}

						return nil
					}); err != nil {
					return err
				}
				return nil
			})
		}
	}
	return g.Wait()
}

func (swgrps SWGResourceProviders) Filter(allowListFile string) (SWGResourceProviders, error) {
	f, err := os.Open(allowListFile)
	if err != nil {
		return SWGResourceProviders{}, err
	}

	newswgrps := SWGResourceProviders{}

	p := regexp.MustCompile(`^([^:]+)$|^([^:]+):([^:]+)$|^([^:]+):([^:]+):([^:]+)$`)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		m := p.FindStringSubmatch(scanner.Text())

		mm := []string{}
		for _, v := range m[1:] {
			if v != "" {
				mm = append(mm, v)
			}
		}

		switch len(mm) {
		case 1:
			rp := mm[0]
			if rp == "*" {
				newswgrps = swgrps.Copy()
			} else {
				newswgrps.CopyRPFrom(swgrps, rp)
			}
		case 2:
			rpName, apiName := mm[0], mm[1]
			if rpName == "*" {
				if apiName == "*" {
					newswgrps = swgrps.Copy()
				} else {
					for rpName := range swgrps {
						newswgrps.CopyAPIFrom(swgrps, rpName, apiName)
					}
				}
			} else {
				if _, ok := swgrps[rpName]; !ok {
					continue
				}
				if apiName == "*" {
					for api := range swgrps[rpName].Apis {
						newswgrps.CopyAPIFrom(swgrps, rpName, api)
					}
				} else {
					newswgrps.CopyAPIFrom(swgrps, rpName, apiName)
				}
			}
		case 3:
			rpName, apiName, schemaName := mm[0], mm[1], mm[2]
			if rpName == "*" {
				if apiName == "*" {
					if schemaName == "*" {
						newswgrps = swgrps.Copy()
					} else {
						for rpName := range swgrps {
							for apiName := range swgrps[rpName].Apis {
								newswgrps.CopySchemaFrom(swgrps, rpName, apiName, schemaName)
							}
						}
					}
				} else {
					for rpName := range swgrps {
						if _, ok := swgrps[rpName].Apis[apiName]; !ok {
							continue
						}

						if schemaName == "*" {
							for schemaName := range swgrps[rpName].Apis[apiName].Schemas {
								newswgrps.CopySchemaFrom(swgrps, rpName, apiName, schemaName)
							}
						} else {
							newswgrps.CopySchemaFrom(swgrps, rpName, apiName, schemaName)
						}
					}
				}
			} else {
				if apiName == "*" {
					if schemaName == "*" {
						newswgrps.CopyRPFrom(swgrps, rpName)
					} else {
						if _, ok := swgrps[rpName]; !ok {
							continue
						}
						for apiName := range swgrps[rpName].Apis {
							newswgrps.CopySchemaFrom(swgrps, rpName, apiName, schemaName)
						}
					}
				} else {
					if schemaName == "*" {
						newswgrps.CopyAPIFrom(swgrps, rpName, apiName)
					} else {
						newswgrps.CopySchemaFrom(swgrps, rpName, apiName, schemaName)
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return SWGResourceProviders{}, err
	}

	return newswgrps, nil
}

func collectAllTFCandidateSchemas(swaggerRepoBaseURI, relPath string) ([]SWGSchema, error) {
	coreSchemas, err := core.CollectSWGSchemas(swaggerRepoBaseURI, relPath, func(swagger *openapispec.Swagger) (schemaNames []string) {
		if swagger.Paths == nil {
			return nil
		}
		schemaNameSet := map[string]struct{}{}
		for _, p := range swagger.Paths.Paths {
			// We only consider resource contains GET, PUT and DELETE methods as a Terraform candidate
			if p.Put == nil || p.Delete == nil || p.Get == nil {
				continue
			}

			for _, param := range p.Put.Parameters {
				if param.In != "body" {
					continue
				}

				// TODO: we should handle cross file reference for these in-body parameters.
				// Here we simply ignore the cross file reference.
				if !param.Schema.Ref.HasFragmentOnly {
					continue
				}

				refString := param.Schema.Ref.GetPointer().String()
				refStringPattern := regexp.MustCompile(`^/definitions/([^/]+)$`)
				matches := refStringPattern.FindStringSubmatch(refString)
				if len(matches) != 2 {
					continue
				}

				schemaNameSet[matches[1]] = struct{}{}
			}
		}
		schemaNames = make([]string, 0, len(schemaNameSet))
		for schemaName := range schemaNameSet {
			schemaNames = append(schemaNames, schemaName)
		}
		return schemaNames
	})

	if err != nil {
		return nil, err
	}

	schemas := make([]SWGSchema, 0, len(coreSchemas))
	for _, coreSchema := range coreSchemas {
		schemas = append(schemas, SWGSchema{coreSchema})
	}
	return schemas, nil
}
