package core

import (
	"fmt"
	"log"
	"strings"

	"github.com/zclconf/go-cty/cty"

	"github.com/magodo/terraform-provider-azurerm-insight/pkg/core/propertyaddr"
)

type SwaggerLink struct {
	Spec       *string                          `json:"swagger,omitempty"` // swagger spec relative path that this propertyaddr resides in, this overrides the global swagger scope
	SchemaProp propertyaddr.SwaggerPropertyAddr `json:"prop"`              // dot-separated swagger schemas propertyaddr, starting from the schemas used as the PUT body parameter
}

type TFSchemaPropertyLinks map[string][]SwaggerLink

type TFSchema struct {
	Name          string
	SwaggerSpec   string `json:"swagger"` // swagger spec relative path path that all the linked swagger property resides in by default
	PropertyLinks TFSchemaPropertyLinks
}

func NewSchema(name string) *TFSchema {
	return &TFSchema{
		Name:          name,
		PropertyLinks: map[string][]SwaggerLink{},
	}
}

func (schema TFSchema) LinkSwagger(swgSchemaCache *SWGSchemas, swaggerBasePath string) error {
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
	if strings.HasPrefix(schema.SwaggerSpec, "/") {
		return fmt.Errorf(`swagger spec path should be relative (not starting with "/")`)
	}
	for tfProp, tfToSwaggerLinks := range schema.PropertyLinks {
		if addr := propertyaddr.NewPropertyAddrFromString(tfProp); addr.Owner() != "" {
			return fmt.Errorf("terraform property addr %s should not specify owner", addr)
		}
		for _, link := range tfToSwaggerLinks {
			if link.SchemaProp.Schema == "" {
				return fmt.Errorf("swagger property addr %s should specify owner", link.SchemaProp)
			}
			if link.Spec != nil && strings.HasPrefix(*link.Spec, "/") {
				return fmt.Errorf(`swagger spec path should be relative (not starting with "/")`)
			}
		}
	}
	return nil
}

// NewSchemaScaffoldFromTerraformBlock construct the TFSchema for a certain resource from the terraform resource block derived
// from `terraform providers schema -json`.
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

// UpdateSchemaScaffoldFromTerraformBlock is similart to NewSchemaScaffoldFromTerraformBlock, except it will update the constrcutred
// TFSchema with the existing TFSchema.
func UpdateSchemaScaffoldFromTerraformBlock(name string, block *TerraformBlock, oldSchema *TFSchema) (*TFSchema, error) {
	newSchema := NewSchemaScaffoldFromTerraformBlock(name, block)
	if oldSchema.Name != name {
		return nil, fmt.Errorf("schema name between existing (%q) and the new (%q) TFSchema is different.", newSchema.Name, oldSchema.Name)
	}

	newSchema.SwaggerSpec = oldSchema.SwaggerSpec

	for propName, propLink := range oldSchema.PropertyLinks {
		if newSchema.PropertyLinks[propName] != nil {
			newSchema.PropertyLinks[propName] = propLink
		} else {
			log.Printf("Warning: The new version of schema %q removed property %q", name, propName)
		}
	}
	return newSchema, nil
}
