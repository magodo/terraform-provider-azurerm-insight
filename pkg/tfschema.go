package pkg

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"

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

func (schema TFSchema) LinkSwagger(swgSchemaCache SWGSchemas, swaggerBasePath string) error {
	for tfProp, tfToSwaggerLinks := range schema.PropertyLinks {
		tfPropAddr := propertyaddr.NewPropertyAddrFromStringWithOwner(schema.Name, tfProp)
		for _, link := range tfToSwaggerLinks {
			swaggerRelPath := schema.SwaggerSpec
			if link.Spec != nil {
				swaggerRelPath = *link.Spec
			}
			// link swgschema
			if err := swgSchemaCache.LinkSWGSchema(swaggerBasePath, swaggerRelPath, link.SchemaProp, *tfPropAddr); err != nil {
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
	for attrKey, attrVal := range block.Attributes {
		addr := parentBlockAddr.Append(attrKey)
		recordAttributeByType(addr, attributes, attrVal.Type)
	}
	for blockKey, blockVal := range block.BlockTypes {
		addr := parentBlockAddr.Append(blockKey)
		recordAttributeWithinBlock(addr, attributes, &blockVal.TerraformBlock)
	}
}

func recordAttributeByType(parentAddr propertyaddr.PropertyAddr, attributes TFSchemaPropertyLinks, elementType *cty.Type) {
	switch {
	case elementType == nil,
		elementType.IsPrimitiveType():
		attributes[parentAddr.String()] = []SwaggerLink{}
	case elementType.IsListType():
		recordAttributeByType(parentAddr, attributes, elementType.ListElementType())
	case elementType.IsSetType():
		recordAttributeByType(parentAddr, attributes, elementType.SetElementType())
	case elementType.IsMapType():
		recordAttributeByType(parentAddr, attributes, elementType.MapElementType())
	case elementType.IsObjectType():
		for attrKey, attrVal := range elementType.AttributeTypes() {
			addr := parentAddr.Append(attrKey)
			recordAttributeByType(addr, attributes, &attrVal)
		}
	}
}
