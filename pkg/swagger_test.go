package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/magodo/terraform-provider-azurerm-insight/pkg/propertyaddr"
	"github.com/stretchr/testify/assert"
)

func TestNewSWGSchema(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	specFooPath := filepath.Join(pwd, "testdata", "swagger", "foo.json")
	specBarPath := filepath.Join(pwd, "testdata", "swagger", "bar.json")
	specFoo, err := LoadSwagger(specFooPath)
	specBar, err := LoadSwagger(specBarPath)
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
				Properties: map[string]*SWGSchemaProperty{
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
				Properties: map[string]*SWGSchemaProperty{
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
				Properties: map[string]*SWGSchemaProperty{
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
				Properties: map[string]*SWGSchemaProperty{
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
				Properties: map[string]*SWGSchemaProperty{
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
				Properties: map[string]*SWGSchemaProperty{
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
				Properties: map[string]*SWGSchemaProperty{
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specBar.Definitions["def_bar"].Properties["prop_primitive"],
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
				Properties: map[string]*SWGSchemaProperty{
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
		{
			specPath:   specFooPath,
			schemaName: "def_allOf",
			err:        nil,
			expect: SWGSchema{
				Name:     "def_allOf",
				SpecPath: specFooPath,
				Properties: map[string]*SWGSchemaProperty{
					"p1": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_allOf"].Properties["p1"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_allOf": struct{}{},
						},
					},
					"prop_nested1": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_allOf"].AllOf[0].Properties["prop_nested1"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_allOf": struct{}{},
						},
					},
					"prop_nested2": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_allOf"].AllOf[0].Properties["prop_nested2"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_allOf": struct{}{},
						},
					},
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specBar.Definitions["def_bar"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_allOf": struct{}{},
							specBarPath + "#/definitions/def_bar":   struct{}{},
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
	specBarPath := filepath.Join(pwd, "testdata", "swagger", "bar.json")
	specFoo, err := LoadSwagger(specFooPath)
	specBar, err := LoadSwagger(specBarPath)
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
				Properties: map[string]*SWGSchemaProperty{
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
				Properties: map[string]*SWGSchemaProperty{
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
				Properties: map[string]*SWGSchemaProperty{
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
				Properties: map[string]*SWGSchemaProperty{
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
			expandAddrs: []propertyaddr.PropertyAddr{
				*propertyaddr.NewPropertyAddrFromStringWithOwner("def_propCrossFileRef", "prop_crossFileRef"),
			},
			err: nil,
			expect: SWGSchema{
				Name:     "def_propCrossFileRef",
				SpecPath: specFooPath,
				Properties: map[string]*SWGSchemaProperty{
					"prop_crossFileRef.prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specBar.Definitions["def_bar"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_propCrossFileRef": struct{}{},
							specBarPath + "#/definitions/def_bar":              struct{}{},
						},
					},
				},
				swagger: specFoo,
			},
		},
		{
			specPath:   specFooPath,
			schemaName: "def_inFileRef",
			expandAddrs: []propertyaddr.PropertyAddr{
				*propertyaddr.NewPropertyAddrFromStringWithOwner("def_inFileRef", "prop_primitive"),
			},
			err: nil,
			expect: SWGSchema{
				Name:     "def_inFileRef",
				SpecPath: specFooPath,
				Properties: map[string]*SWGSchemaProperty{
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
			expandAddrs: []propertyaddr.PropertyAddr{
				*propertyaddr.NewPropertyAddrFromStringWithOwner("def_crossFileRef", "prop_primitive"),
			},
			err: nil,
			expect: SWGSchema{
				Name:     "def_crossFileRef",
				SpecPath: specFooPath,
				Properties: map[string]*SWGSchemaProperty{
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specBar.Definitions["def_bar"].Properties["prop_primitive"],
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
			expandAddrs: []propertyaddr.PropertyAddr{
				*propertyaddr.NewPropertyAddrFromStringWithOwner("def_selfRef", ""),
			},
			expect: SWGSchema{
				Name:     "def_selfRef",
				SpecPath: specFooPath,
				Properties: map[string]*SWGSchemaProperty{
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
		swgschema, err := NewSWGSchema(c.specPath, c.schemaName)
		assert.Equal(t, c.err, err, idx)
		if err == nil {
			for _, addr := range c.expandAddrs {
				assert.NoError(t, swgschema.ExpandPropertyOneLevelDeep(addr), idx)
			}
			assert.Equal(t, c.expect, *swgschema, idx)
		}
	}
}

func TestLinkSWGSchema_AddTFLink(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	specFooPath := filepath.Join(pwd, "testdata", "swagger", "foo.json")
	specBarPath := filepath.Join(pwd, "testdata", "swagger", "bar.json")
	specFoo, err := LoadSwagger(specFooPath)
	specBar, err := LoadSwagger(specBarPath)
	_ = specBar
	if err != nil {
		t.Fatal(err)
	}

	type step struct {
		swgPropAddr propertyaddr.PropertyAddr
		tfPropAddr  propertyaddr.PropertyAddr
		err         error
		expect      SWGSchema
	}

	cases := []struct {
		specPath   string
		schemaName string
		steps      []step
	}{
		{
			specPath:   specFooPath,
			schemaName: "def_a",
			steps: []step{
				{
					swgPropAddr: *propertyaddr.NewPropertyAddrFromString("def_a:prop_primitive"),
					tfPropAddr:  *propertyaddr.NewPropertyAddrFromString("res1.p1"),
					err:         nil,
					expect: SWGSchema{
						Name:     "def_a",
						SpecPath: specFooPath,
						Properties: map[string]*SWGSchemaProperty{
							"prop_primitive": {
								TFLinks: []TFLink{
									{*propertyaddr.NewPropertyAddrFromString("res1.p1")},
								},
								schema: specFoo.Definitions["def_foo"].Properties["prop_primitive"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_a":   struct{}{},
									specFooPath + "#/definitions/def_foo": struct{}{},
								},
							},
							"p1": {
								TFLinks: []TFLink{},
								schema:  specFoo.Definitions["def_a"].Properties["p1"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_a": struct{}{},
								},
							},
						},
						swagger: specFoo,
					},
				},
				// Add a second tf link to the same swg property
				{
					swgPropAddr: *propertyaddr.NewPropertyAddrFromString("def_a:prop_primitive"),
					tfPropAddr:  *propertyaddr.NewPropertyAddrFromString("res2.p1"),
					err:         nil,
					expect: SWGSchema{
						Name:     "def_a",
						SpecPath: specFooPath,
						Properties: map[string]*SWGSchemaProperty{
							"prop_primitive": {
								TFLinks: []TFLink{
									{*propertyaddr.NewPropertyAddrFromString("res1.p1")},
									{*propertyaddr.NewPropertyAddrFromString("res2.p1")},
								},
								schema: specFoo.Definitions["def_foo"].Properties["prop_primitive"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_a":   struct{}{},
									specFooPath + "#/definitions/def_foo": struct{}{},
								},
							},
							"p1": {
								TFLinks: []TFLink{},
								schema:  specFoo.Definitions["def_a"].Properties["p1"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_a": struct{}{},
								},
							},
						},
						swagger: specFoo,
					},
				},
				{
					swgPropAddr: *propertyaddr.NewPropertyAddrFromString("def_a:p1.prop_primitive"),
					tfPropAddr:  *propertyaddr.NewPropertyAddrFromString("res1.p2"),
					err:         nil,
					expect: SWGSchema{
						Name:     "def_a",
						SpecPath: specFooPath,
						Properties: map[string]*SWGSchemaProperty{
							"prop_primitive": {
								TFLinks: []TFLink{
									{*propertyaddr.NewPropertyAddrFromString("res1.p1")},
									{*propertyaddr.NewPropertyAddrFromString("res2.p1")},
								},
								schema: specFoo.Definitions["def_foo"].Properties["prop_primitive"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_a":   struct{}{},
									specFooPath + "#/definitions/def_foo": struct{}{},
								},
							},
							"p1.prop_primitive": {
								TFLinks: []TFLink{
									{*propertyaddr.NewPropertyAddrFromString("res1.p2")},
								},
								schema: specFoo.Definitions["def_foo"].Properties["prop_primitive"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_a":   struct{}{},
									specBarPath + "#/definitions/def_bar": struct{}{},
								},
							},
							"p1.p1_1": {
								TFLinks: []TFLink{},
								schema:  specFoo.Definitions["def_foo"].Properties["p1"].Properties["p1_1"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_a": struct{}{},
								},
							},
						},
						swagger: specFoo,
					},
				},
			},
		},
	}

	for idx, c := range cases {
		schema, err := NewSWGSchema(c.specPath, c.schemaName)
		assert.NoError(t, err, idx)
		if err == nil {
			for iidx, s := range c.steps {
				err := schema.AddTFLink(s.swgPropAddr, s.tfPropAddr)
				assert.Equal(t, s.err, err, fmt.Sprintf("%d.%d", idx, iidx))
				if err == nil {
					assert.Equal(t, s.expect, *schema, fmt.Sprintf("%d.%d", idx, iidx))
				}
			}
		}
	}
}
