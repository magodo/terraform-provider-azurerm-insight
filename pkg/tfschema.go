package pkg

import (
	"fmt"
)

type SwaggerLink struct {
	Spec       *string      `json:"spec"` // swagger spec alias that this property resides in, this overrides the global swagger spec scope
	SchemaProp propertyAddr `json:"prop"` // dot-separated swagger schema property, starting from the schema used as the PUT body parameter
}

type TFSchemaPropertyLinks map[string][]SwaggerLink

type TFSchema struct {
	Name          string
	SwaggerSpec   string `json:"spec"`
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
		tfPropAddr := newPropertyAddrFromString(tfProp)
		tfPropAddr.owner = schema.Name
		for _, link := range tfToSwaggerLinks {
			specPath := schema.SwaggerSpec
			if link.Spec != nil {
				specPath = *link.Spec
			}
			// valid property links will always have the form: <swagger schema definition>.<prop1>.<prop2>...
			switch len(link.SchemaProp.addrs) {
			case 0:
				return fmt.Errorf("empty property link found for %s, please remove it", schema.Name)
			case 1:
				prop := link.SchemaProp.addrs[0]
				return fmt.Errorf("malformed property link found for %s: %s", schema.Name, prop)
			}

			swgSchemaName, swgProps := link.SchemaProp.addrs[0], link.SchemaProp.addrs[1:]
			swagPropAddr := newPropertyAddr(swgSchemaName, swgProps...)

			// link swgschema
			if err := LinkSWGSchema(specPath, *swagPropAddr, *tfPropAddr); err != nil {
				return fmt.Errorf("linking swgschema: %w", err)
			}
		}
	}

	return nil
}

func NewSchemaScaffoldFromTerraformBlock(name string, block *TerraformBlock) *TFSchema {
	schema := NewSchema(name)
	recordAttributeWithinBlock(propertyAddr{}, schema.PropertyLinks, block)
	return schema
}

func recordAttributeWithinBlock(parentBlockAddr propertyAddr, attributes TFSchemaPropertyLinks, block *TerraformBlock) {
	for attrKey := range block.Attributes {
		addr := parentBlockAddr.Append(attrKey)
		attributes[addr.String()] = []SwaggerLink{}
	}
	for blockKey, blockVal := range block.BlockTypes {
		addr := parentBlockAddr.Append(blockKey)
		recordAttributeWithinBlock(addr, attributes, &blockVal.TerraformBlock)
	}
}
