package pkg

import (
	"github.com/magodo/terraform-provider-azurerm-insight/pkg/propertyaddr"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestNewSWGSchema(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	specFooPath := filepath.Join(pwd, "testdata", "swagger", "foo.json")
	specBarPath := filepath.Join(pwd, "testdata", "swagger", "bar.json")
	specFoo, err := LoadSwagger(specFooPath)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		specPath   string
		schemaName string
		err        error
		expect     SWGSchema
	}{
		{
			specPath:   specFooPath,
			schemaName: "def_foo",
			err:        nil,
			expect: SWGSchema{
				Name:     "def_foo",
				SpecPath: specFooPath,
				Properties: map[string]SWGSchemaProperty{
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_foo"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_foo": struct{}{},
						},
					},
				},
				swagger: specFoo,
			},
		},
		{
			specPath:   specFooPath,
			schemaName: "def_regular",
			err:        nil,
			expect: SWGSchema{
				Name:     "def_regular",
				SpecPath: specFooPath,
				Properties: map[string]SWGSchemaProperty{
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_regular"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_regular": struct{}{},
						},
					},
					"prop_array_of_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_regular"].Properties["prop_array_of_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_regular": struct{}{},
						},
					},
					"prop_object": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_regular"].Properties["prop_object"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_regular": struct{}{},
						},
					},
				},
				swagger: specFoo,
			},
		},
		{
			specPath:   specFooPath,
			schemaName: "def_propInFileRef",
			err:        nil,
			expect: SWGSchema{
				Name:     "def_propInFileRef",
				SpecPath: specFooPath,
				Properties: map[string]SWGSchemaProperty{
					"prop_inFileRef": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_propInFileRef"].Properties["prop_inFileRef"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_propInFileRef": struct{}{},
						},
					},
				},
				swagger: specFoo,
			},
		},
		{
			specPath:   specFooPath,
			schemaName: "def_propSelfRef",
			err:        nil,
			expect: SWGSchema{
				Name:     "def_propSelfRef",
				SpecPath: specFooPath,
				Properties: map[string]SWGSchemaProperty{
					"prop_selfRef": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_propSelfRef"].Properties["prop_selfRef"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_propSelfRef": struct{}{},
						},
					},
				},
				swagger: specFoo,
			},
		},
		{
			specPath:   specFooPath,
			schemaName: "def_propCrossFileRef",
			err:        nil,
			expect: SWGSchema{
				Name:     "def_propCrossFileRef",
				SpecPath: specFooPath,
				Properties: map[string]SWGSchemaProperty{
					"prop_crossFileRef": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_propCrossFileRef"].Properties["prop_crossFileRef"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_propCrossFileRef": struct{}{},
						},
					},
				},
				swagger: specFoo,
			},
		},
		{
			specPath:   specFooPath,
			schemaName: "def_inFileRef",
			err:        nil,
			expect: SWGSchema{
				Name:     "def_inFileRef",
				SpecPath: specFooPath,
				Properties: map[string]SWGSchemaProperty{
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_inFileRef"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_inFileRef": struct{}{},
							specFooPath + "#/definitions/def_foo":       struct{}{},
						},
					},
				},
				swagger: specFoo,
			},
		},
		{
			specPath:   specFooPath,
			schemaName: "def_crossFileRef",
			err:        nil,
			expect: SWGSchema{
				Name:     "def_crossFileRef",
				SpecPath: specFooPath,
				Properties: map[string]SWGSchemaProperty{
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_crossFileRef"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_crossFileRef": struct{}{},
							specBarPath + "#/definitions/def_bar":          struct{}{},
						},
					},
				},
				swagger: specFoo,
			},
		},
		{
			specPath:   specFooPath,
			schemaName: "def_selfRef",
			err:        nil,
			expect: SWGSchema{
				Name:     "def_selfRef",
				SpecPath: specFooPath,
				Properties: map[string]SWGSchemaProperty{
					"": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_selfRef"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_selfRef": struct{}{},
						},
					},
				},
				swagger: specFoo,
			},
		},
	}

	for idx, c := range cases {
		actual, err := NewSWGSchema(c.specPath, c.schemaName)
		assert.Equal(t, c.err, err, idx)
		if err == nil {
			assert.Equal(t, c.expect, *actual, idx)
		}
	}
}

func TestSWGSchema_ExpandPropertyOneLevelDeep(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	specFooPath := filepath.Join(pwd, "testdata", "swagger", "foo.json")
	//specBarPath := filepath.Join(pwd, "testdata", "swagger", "bar.json")
	specFoo, err := LoadSwagger(specFooPath)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		specPath    string
		schemaName  string
		expandAddrs []propertyaddr.PropertyAddr
		err         error
		expect      SWGSchema
	}{
		{
			specPath:   specFooPath,
			schemaName: "def_foo",
			expandAddrs: []propertyaddr.PropertyAddr{
				*propertyaddr.NewPropertyAddrFromStringWithOwner("def_foo", "prop_primitive"),
				*propertyaddr.NewPropertyAddrFromStringWithOwner("def_foo", "prop_primitive"),
			},
			err: nil,
			expect: SWGSchema{
				Name:     "def_foo",
				SpecPath: specFooPath,
				Properties: map[string]SWGSchemaProperty{
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_foo"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_foo": struct{}{},
						},
					},
				},
				swagger: specFoo,
			},
		},
		{
			specPath:   specFooPath,
			schemaName: "def_regular",
			expandAddrs: []propertyaddr.PropertyAddr{
				*propertyaddr.NewPropertyAddrFromStringWithOwner("def_regular", "prop_object"),
			},
			err: nil,
			expect: SWGSchema{
				Name:     "def_regular",
				SpecPath: specFooPath,
				Properties: map[string]SWGSchemaProperty{
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_regular"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_regular": struct{}{},
						},
					},
					"prop_array_of_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_regular"].Properties["prop_array_of_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_regular": struct{}{},
						},
					},
					"prop_object.prop_nested": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_regular"].Properties["prop_object"].Properties["prop_nested"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_regular": struct{}{},
						},
					},
				},
				swagger: specFoo,
			},
		},
		{
			specPath:   specFooPath,
			schemaName: "def_propInFileRef",
			expandAddrs: []propertyaddr.PropertyAddr{
				*propertyaddr.NewPropertyAddrFromStringWithOwner("def_propInFileRef", "prop_inFileRef"),
			},
			err: nil,
			expect: SWGSchema{
				Name:     "def_propInFileRef",
				SpecPath: specFooPath,
				Properties: map[string]SWGSchemaProperty{
					"prop_inFileRef.prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_foo"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_propInFileRef": struct{}{},
							specFooPath + "#/definitions/def_foo":           struct{}{},
						},
					},
				},
				swagger: specFoo,
			},
		},
		{
			specPath:   specFooPath,
			schemaName: "def_propSelfRef",
			expandAddrs: []propertyaddr.PropertyAddr{
				*propertyaddr.NewPropertyAddrFromStringWithOwner("def_propSelfRef", "prop_selfRef"),
			},
			err: nil,
			expect: SWGSchema{
				Name:     "def_propSelfRef",
				SpecPath: specFooPath,
				Properties: map[string]SWGSchemaProperty{
					"prop_selfRef": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_propSelfRef"].Properties["prop_selfRef"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_propSelfRef": struct{}{},
						},
					},
				},
				swagger: specFoo,
			},
		},
	}

	for idx, c := range cases {
		swgschema, err := NewSWGSchema(c.specPath, c.schemaName)
		assert.Equal(t, c.err, err, idx)
		if err == nil {
			for _, addr := range c.expandAddrs {
				assert.NoError(t, swgschema.ExpandPropertyOneLevelDeep(addr), idx)
			}
			assert.Equal(t, c.expect, *swgschema)
		}
	}
}
