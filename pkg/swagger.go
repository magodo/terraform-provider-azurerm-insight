package pkg

import (
	"fmt"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"sync"
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

func LinkSWGSchema(spec, def string, swgprop propertyAddr, tfresource string) error {
	swgSchemaCache.Lock()
	defer swgSchemaCache.Unlock()

	// construct key
	k := spec + "-" + def
	swgSchema, ok := swgSchemaCache.m[k]
	if !ok {
		// load the schema definition from swagger spec and initialize the swgschema based on it
		swaggerSchema, err := LoadSwaggerSchema(spec, def)
		if err != nil {
			return fmt.Errorf("loading swagger schema definition: %w", err)
		}

		swgSchema = &SWGSchema{
			Name:          def,
			SpecPath:      spec,
			PropertyLinks: map[string][]TFLink{},
		}

		// we only set the first-level properties (including inherited) on initialization
		topProps, err := swaggerSchemaTopProperties(spec, swaggerSchema)
		if err != nil {
			return fmt.Errorf("getting top level properties for %s (%s)", def, spec)
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
func swaggerSchemaTopProperties(base string, schema *spec.Schema) (map[string]spec.Schema, error) {
	out := map[string]spec.Schema{}
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

type SwaggerSchemaCache struct {
	sync.Mutex
	m map[string]*spec.Schema
}

var swaggerSchemaCache SwaggerSchemaCache

// LoadSwaggerSchema load a certain swagger schema definition
func LoadSwaggerSchema(spec, def string) (*spec.Schema, error) {
	swaggerSchemaCache.Lock()
	defer swaggerSchemaCache.Unlock()

	// construct key
	k := spec + "-" + def
	if schema, ok := swaggerSchemaCache.m[k]; ok {
		return schema, nil
	}

	// load the schema from swagger spec and cache it
	specDoc, err := loads.Spec(spec)
	if err != nil {
		return nil, fmt.Errorf("loading swagger spec %s: %w", spec, err)
	}

	schema := specDoc.Schema()
	swaggerSchemaCache.m[k] = schema
	return schema, nil
}
