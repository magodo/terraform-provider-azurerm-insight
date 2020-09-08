package pkg

import (
	"fmt"
	"sync"

	"github.com/go-openapi/loads"
	openapispec "github.com/go-openapi/spec"
)

type TFLink struct {
	Resource string
	Prop     propertyAddr
}

type SWGSchemaPropertyLinks map[string][]TFLink

type SWGSchema struct {
	Name          string
	SpecPath      string
	PropertyLinks SWGSchemaPropertyLinks
}

type SWGSchemaCache struct {
	sync.Mutex
	m map[string]*SWGSchema
}

var swgSchemaCache SWGSchemaCache

func LinkSWGSchema(specPath string, swgPropAddr, tfPropAddr propertyAddr) error {
	if !swgPropAddr.IsCanonical() {
		return fmt.Errorf("swagger property address is not canonical: %s", swgPropAddr.String())
	}

	if !tfPropAddr.IsCanonical() {
		return fmt.Errorf("terraform property address is not canonical: %s", tfPropAddr.String())
	}

	swgSchemaCache.Lock()
	defer swgSchemaCache.Unlock()

	// construct key
	k := specPath + "-" + swgPropAddr.owner
	swgSchema, ok := swgSchemaCache.m[k]
	if !ok {
		// load the swagger specPath and initialize the swgschema based on it
		spec, err := LoadSwaggerSpec(specPath)
		if err != nil {
			return fmt.Errorf("loading swagger schema definition: %w", err)
		}

		swgSchema = &SWGSchema{
			Name:          swgPropAddr.owner,
			SpecPath:      specPath,
			PropertyLinks: map[string][]TFLink{},
		}

		swaggerRef, err := swgPropAddr.ToDefinitionRef()
		if err != nil {
			return fmt.Errorf("construct definition reference for %s: %w", swgPropAddr.String(), err)
		}
		swaggerSchema, err := openapispec.ResolveRefWithBase(spec, &swaggerRef, &openapispec.ExpandOptions{RelativeBase: specPath})
		if err != nil {
			return fmt.Errorf("resolve reference %s: %w", swaggerRef.String(), err)
		}

		// we only set the first-level properties (including inherited) on initialization
		topProps, err := swaggerSchemaTopProperties(specPath, swaggerSchema)
		if err != nil {
			return fmt.Errorf("getting top level properties for %s (%s)", swgPropAddr.String(), specPath)
		}
		for prop := range topProps {
			swgSchema.PropertyLinks[prop] = []TFLink{}
		}
	}

	// TODO: now we have the swgSchema, we need to link the terraform resource's properties to it

	// expand to the level of swgprop

	swgSchemaCache.m[k] = swgSchema

	return nil
}

// swaggerSchemaTopProperties get the top level properties of the input swagger schema.
// It expands both the inherited schema or the refs, but only for the top level.
func swaggerSchemaTopProperties(base string, schema *openapispec.Schema) (map[string]openapispec.Schema, error) {
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

		// TODO: follow the ref and take its top properties
		// NOTE: consider cross file reference and cyclic reference
	}
	return out, nil
}

type SwaggerSpecCache struct {
	sync.Mutex
	m map[string]*loads.Document
}

var swaggerSpecCache SwaggerSpecCache

// LoadSwaggerSpec load a certain swagger spec (document)
func LoadSwaggerSpec(spec string) (*loads.Document, error) {
	swaggerSpecCache.Lock()
	defer swaggerSpecCache.Unlock()

	// construct key
	if schema, ok := swaggerSpecCache.m[spec]; ok {
		return schema, nil
	}

	doc, err := loads.Spec(spec)
	if err != nil {
		return nil, fmt.Errorf("loading swagger spec %s: %w", spec, err)
	}

	swaggerSpecCache.m[spec] = doc
	return doc, nil
}
