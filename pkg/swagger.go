package pkg

import (
	"fmt"
	"github.com/go-openapi/loads"
	openapispec "github.com/go-openapi/spec"
	"github.com/magodo/terraform-provider-azurerm-insight/pkg/propertyaddr"
	"sync"
)

type TFLink struct {
	Prop propertyaddr.PropertyAddr
}

type SWGSchemaPropertyLink struct {
	// Terraform property addresses
	TFLink []TFLink
	
	// The schema of this swagger schema property
	schema openapispec.Schema
	
	// The resolved URI refs along the way to this schema, each is an absolute/normalized reference.
	resolvedRefs []string
}

// (key: swagger schema relative property addr)
type SWGSchemaPropertyLinks map[string]SWGSchemaPropertyLink

type SWGSchema struct {
	Name          string
	SpecPath      string
	PropertyLinks SWGSchemaPropertyLinks

	spec *loads.Document
}

// ExpandPropertyLinkOneLevel expand the specified swagger schema property one level, with any allOf and $ref taken into consideration.
// If `addr` is absent in the SWGSchema's property links map, error will be returned.
func (s *SWGSchema) ExpandPropertyLinkOneLevel(parentAddr propertyaddr.PropertyAddr) error {
	raddr := parentAddr.RelativeAddrs()
	link, ok := s.PropertyLinks[raddr.String()]
	if !ok {
		return fmt.Errorf("swagger schema %s (spec: %s) doesn't have property %s", s.Name, s.SpecPath, raddr)
	}

	// Direct properties
	for propK, propV := range link.schema.Properties {
		addr := parentAddr.Append(propK)	
		s.PropertyLinks[addr] = SWGSchemaPropertyLink{
			TFLink:       []string{link.TFLink...},
			schema:       openapispec.Schema{},
			resolvedRefs: nil,
		}
	}
	
	// TODO:
	// 1. keep the property link of the expanded property in the expanded-out properties
	// 2. keep the resolvedRefs of the expanded property in the expanded-out properties
	// 3. remove the expanded property from PropertyLinks
	// 3. remove the expanded property from resolvedRefs
	_ = links

	return nil
}

type SWGSpecSchemaCache struct {
	sync.Mutex
	m map[string]*SWGSchema
}

func (c *SWGSpecSchemaCache) Lock() {
	c.Mutex.Lock()
}

func (c *SWGSpecSchemaCache) Unlock() {
	c.Mutex.Unlock()
}

func (c *SWGSpecSchemaCache) Get(specPath, schemaName string) *SWGSchema {
	k := specPath + "-" + schemaName
	return c.m[k]
}

func (c *SWGSpecSchemaCache) Set(specPath, schemaName string, schema *SWGSchema) {
	k := specPath + "-" + schemaName
	c.m[k] = schema
}

// swgSpecSchemaCache caches the SWGSchema using spec + schema as key.
// During each link operation from terraform schema to swagger schema, it will manipulate one of
// the SWGSchema in this cache. Afterwards, this cache contains all the mapping info from swagger to terraform.
var swgSpecSchemaCache SWGSpecSchemaCache

func LinkSWGSchema(specPath string, swgPropAddr, tfPropAddr propertyaddr.PropertyAddr) error {
	swgSpecSchemaCache.Lock()
	defer swgSpecSchemaCache.Unlock()

	// construct key
	swgSchema := swgSpecSchemaCache.Get(specPath, swgPropAddr.Owner())
	if swgSchema == nil {
		spec, err := LoadSwaggerSpec(specPath)
		if err != nil {
			return err
		}

		swgSchema = &SWGSchema{
			Name:          swgPropAddr.Owner(),
			SpecPath:      specPath,
			PropertyLinks: map[string]SWGSchemaPropertyLink{},
			spec:          spec,
		}

		// expand the root level properties

	}

	// TODO: now we have the swgSchema, we need to link the terraform resource's properties to it

	swgSpecSchemaCache.Set(specPath, swgPropAddr.Owner(), swgSchema)

	return nil
}

// swaggerSchemaTopProperties get the top level properties of the input swagger schema.
// It expands both the inherited schema or the refs, but only for the top level.
func swaggerSchemaTopProperties(specPath string, schema *openapispec.Schema, resolvedRef map[string][]string) (map[string]openapispec.Schema, error) {
	out := map[string]openapispec.Schema{}
	for propK, propV := range schema.Properties {
		out[propK] = propV
	}
	for _, inherit := range schema.AllOf {
		if inherit.Ref.String() == "" {
			for propK, propV := range inherit.Properties {
				out[propK] = propV
			}
			continue
		}

		// follow the ref and take its top properties
		inheritExpandSchema, err := openapispec.ResolveRefWithBase(root, &inherit.Ref, &openapispec.ExpandOptions{RelativeBase: specPath})
		if err != nil {
			return nil, fmt.Errorf("resolve reference %s: %w", inherit.Ref.String(), err)
		}
		normalizeFileRef(&inherit.Ref, specPath).String()
		for propK, propV := range inheritExpandSchema.Properties {
		}
		// NOTE: consider cross file reference and cyclic reference
	}
	return out, nil
}
