package pkg

import (
	"fmt"
	"path/filepath"

	"github.com/magodo/terraform-provider-azurerm-insight/pkg/propertyaddr"
)

type SwaggerLink struct {
	Spec       *string                   `json:"swagger"` // swagger spec abs path that this propertyaddr resides in, this overrides the global swagger scope
	SchemaProp propertyaddr.PropertyAddr `json:"prop"`    // dot-separated swagger schemas propertyaddr, starting from the schemas used as the PUT body parameter
}

type TFSchemaPropertyLinks map[string][]SwaggerLink

type TFSchema struct {
	Name          string
	SwaggerSpec   string `json:"swagger"` // swagger spec abs path that all the linked swagger property resides in by default
	PropertyLinks TFSchemaPropertyLinks
}

func NewSchema(name string) *TFSchema {
	return &TFSchema{
		Name:          name,
		PropertyLinks: map[string][]SwaggerLink{},
	}
}

func (schema TFSchema) LinkSwagger(swaggerBasePath string) error {
	for tfProp, tfToSwaggerLinks := range schema.PropertyLinks {
		tfPropAddr := propertyaddr.NewPropertyAddrFromStringWithOwner(schema.Name, tfProp)
		for _, link := range tfToSwaggerLinks {
			specPath := schema.SwaggerSpec
			if link.Spec != nil {
				specPath = *link.Spec
			}
			specPath = filepath.Join(swaggerBasePath, specPath)
			// link swgschema
			if err := LinkSWGSchema(specPath, link.SchemaProp, *tfPropAddr); err != nil {
				return fmt.Errorf("linking swgschema: %w", err)
			}
		}
	}

	return nil
}

// Validate validates the swagger property and tf schemas property has the correct form
func (schema TFSchema) Validate() error {
	for tfProp, tfToSwaggerLinks := range schema.PropertyLinks {
		if addr := propertyaddr.NewPropertyAddrFromString(tfProp); addr.Owner() != "" {
			return fmt.Errorf("terraform property addr %s should not specify owner", addr)
		}
		for _, link := range tfToSwaggerLinks {
			if link.SchemaProp.Owner() == "" {
				return fmt.Errorf("swagger property addr %s should specify owner", link.SchemaProp)
			}
		}
	}
	return nil
}

func NewSchemaScaffoldFromTerraformBlock(name string, block *TerraformBlock) *TFSchema {
	schema := NewSchema(name)
	recordAttributeWithinBlock(propertyaddr.PropertyAddr{}, schema.PropertyLinks, block)
	return schema
}

func recordAttributeWithinBlock(parentBlockAddr propertyaddr.PropertyAddr, attributes TFSchemaPropertyLinks, block *TerraformBlock) {
	for attrKey := range block.Attributes {
		addr := parentBlockAddr.Append(attrKey)
		attributes[addr.String()] = []SwaggerLink{}
	}
	for blockKey, blockVal := range block.BlockTypes {
		addr := parentBlockAddr.Append(blockKey)
		recordAttributeWithinBlock(addr, attributes, &blockVal.TerraformBlock)
	}
}
