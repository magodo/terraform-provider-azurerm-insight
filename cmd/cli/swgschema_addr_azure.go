package main

import (
	"github.com/magodo/terraform-provider-azurerm-insight/pkg/core"
	"strings"
)

type SWGResourceProviders map[string]SWGResourceProvider

type SWGResourceProvider struct {
	Apis SWGResourceProviderAPIs
}

type SWGResourceProviderAPIs map[string]SWGResourceProviderAPI

type SWGResourceProviderAPI struct {
	Schemas SWGSchemas
}

type SWGSchemas map[string]SWGSchema

type SWGSchema struct {
	core.SWGSchema
}

type SWGSchemaAddr struct {
	raw core.SWGSchemaAddr
	ResourceProvider string
	ApiVersion string
	ResourceCollectionName string
	SchemaName string
}

func NewSWGSchemaAddr(addr core.SWGSchemaAddr) SWGSchemaAddr {
	swaggerRelPath, schemaName := addr.SwaggerRelPath(), addr.SchemaName()
	segments := strings.Split(swaggerRelPath, "/")
	rp, resourceCollectionName, api := segments[0], strings.TrimSuffix(segments[len(segments)-1], ".json"), segments[len(segments)-2]
	return SWGSchemaAddr{
		raw:              addr,
		ResourceProvider: rp,
		ApiVersion:       api,
		ResourceCollectionName:  resourceCollectionName,
		SchemaName:       schemaName,
	}
}

func NewSWGResourceProviders(swgschemas core.SWGSchemas) SWGResourceProviders {
	out := map[string]SWGResourceProvider{}
	for addr, swgschema := range swgschemas.GetAll() {
		addr := NewSWGSchemaAddr(addr)

		if _, ok := out[addr.ResourceProvider]; !ok {
			out[addr.ResourceProvider] = SWGResourceProvider{Apis: map[string]SWGResourceProviderAPI{}}
		}

		if _, ok := out[addr.ResourceProvider].Apis[addr.ApiVersion]; !ok {
			out[addr.ResourceProvider].Apis[addr.ApiVersion] = SWGResourceProviderAPI{Schemas: map[string]SWGSchema{}}
		}


		out[addr.ResourceProvider].Apis[addr.ApiVersion].Schemas[addr.SchemaName] = SWGSchema{*swgschema}
	}

	return out
}
