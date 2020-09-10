package pkg

import (
	"fmt"
	"github.com/magodo/terraform-provider-azurerm-insight/pkg/propertyaddr"
)

type SwaggerLink struct {
	Spec       *string                   `json:"swagger"` // swagger swagger alias that this propertyaddr resides in, this overrides the global swagger swagger scope
	SchemaProp propertyaddr.PropertyAddr `json:"prop"` // dot-separated swagger schema propertyaddr, starting from the schema used as the PUT body parameter
}

type TFSchemaPropertyLinks map[string][]SwaggerLink

type TFSchema struct {
	Name          string
	SwaggerSpec   string `json:"swagger"`
	PropertyLinks TFSchemaPropertyLinks
}

func NewSchema(name string) *TFSchema {
	return &TFSchema{
		Name:          name,
		PropertyLinks: map[string][]SwaggerLink{},
	}
}

func (schema TFSchema) LinkSwagger() error {
	for tfProp, tfToSwaggerLinks := range schema.PropertyLinks {
		tfPropAddr := propertyaddr.NewPropertyAddrFromStringWithOwner(schema.Name, tfProp)
		for _, link := range tfToSwaggerLinks {
			specPath := schema.SwaggerSpec
			if link.Spec != nil {
				specPath = *link.Spec
			}
			// link swgschema
			if err := LinkSWGSchema(specPath, link.SchemaProp, *tfPropAddr); err != nil {
				return fmt.Errorf("linking swgschema: %w", err)
			}
		}
	}

	return nil
}

func (schema TFSchema) Validate() error {
	if err := specAlias.ValidateAlias(schema.SwaggerSpec); err != nil {
		return err
	}
	for tfProp, tfToSwaggerLinks := range schema.PropertyLinks {
		if addr := propertyaddr.NewPropertyAddrFromString(tfProp); addr.Owner() != "" {
			return fmt.Errorf("terraform property addr %s should not specify owner", addr)
		}
		for _, link := range tfToSwaggerLinks {
			if link.Spec != nil {
				if err := specAlias.ValidateAlias(*link.Spec); err != nil {
					return err
				}
			}
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
