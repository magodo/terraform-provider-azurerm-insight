package pkg

import (
	"fmt"
	"strings"
)

const addrSep = "."

type SwaggerLink struct {
	Spec              *string      `json:"spec"` // swagger spec alias that this property resides in, this overrides the global swagger spec scope
	SwaggerSchemaProp propertyAddr `json:"prop"` // dot-separated property, starting from the schema used as the PUT body parameter
}

type propertyAddr struct {
	segments []string
}

func newPropertyAddrFromString(addr string) *propertyAddr {
	return &propertyAddr{strings.Split(addr, addrSep)}
}

func newPropertyAddr(segments []string) *propertyAddr {
	return &propertyAddr{segments}
}

func (addr propertyAddr) String() string {
	return strings.Join(addr.segments, addrSep)
}

func (addr *propertyAddr) UnmarshalJSON(b []byte) error {
	addr.segments = strings.Split(string(b), addrSep)
	return nil
}

func (addr propertyAddr) MarshalJSON() ([]byte, error) {
	return []byte(strings.Join(addr.segments, addrSep)), nil
}

func (addr propertyAddr) Append(oaddr string) propertyAddr {
	segments := make([]string, len(addr.segments))
	copy(segments, addr.segments)
	segments = append(segments, oaddr)
	return propertyAddr{segments: segments}
}

func (addr *propertyAddr) Pop() string {
	if len(addr.segments) > 1 {
		addr.segments = addr.segments[:len(addr.segments)-1]
	}
	return addr.segments[len(addr.segments)-1]
}

type SwaggerSchemaProp struct {
	segments []string
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
	for tfPropAddr, tfToSwaggerLinks := range schema.PropertyLinks {
		for _, link := range tfToSwaggerLinks {
			specPath := schema.SwaggerSpec
			if link.Spec != nil {
				specPath = *link.Spec
			}
			// valid property links will always have the form: <swagger schema definition>.<prop1>.<prop2>...
			switch len(link.SwaggerSchemaProp.segments) {
			case 0:
				return fmt.Errorf("empty property link found for %s, please remove it", schema.Name)
			case 1:
				prop := link.SwaggerSchemaProp.segments[0]
				return fmt.Errorf("malformed property link found for %s: %s", schema.Name, prop)
			}

			schemaName, swgprops := link.SwaggerSchemaProp.segments[0], newPropertyAddr(link.SwaggerSchemaProp.segments[1:])

			// link swgschema
			if err := LinkSWGSchema(specPath, schemaName, *swgprops, schema.Name); err != nil {
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
