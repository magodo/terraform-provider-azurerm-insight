package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/magodo/terraform-provider-azurerm-insight/pkg/propertyaddr"
)

func TestNewSWGSchema(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	specFooPath := filepath.Join(pwd, "testdata", "swagger", "foo.json")
	specBarPath := filepath.Join(pwd, "testdata", "swagger", "bar.json")
	specFoo, err := LoadSwagger(specFooPath)
	require.NoError(t, err)
	specBar, err := LoadSwagger(specBarPath)
	require.NoError(t, err)

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
					"prop_array_of_object": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_regular"].Properties["prop_array_of_object"],
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
		{
			specPath:   specFooPath,
			schemaName: "def_array_simple",
			err:        nil,
			expect: SWGSchema{
				Name:     "def_array_simple",
				SpecPath: specFooPath,
				Properties: map[string]*SWGSchemaProperty{
					"": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_array_simple"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_array_simple": struct{}{},
						},
					},
				},
				swagger: specFoo,
			},
		},
		{
			specPath:   specFooPath,
			schemaName: "def_array_ref",
			err:        nil,
			expect: SWGSchema{
				Name:     "def_array_ref",
				SpecPath: specFooPath,
				Properties: map[string]*SWGSchemaProperty{
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_foo"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_foo":       struct{}{},
							specFooPath + "#/definitions/def_array_ref": struct{}{},
						},
					},
				},
				swagger: specFoo,
			},
		},
		{
			specPath:   specFooPath,
			schemaName: "def_array_ref_ref",
			err:        nil,
			expect: SWGSchema{
				Name:     "def_array_ref_ref",
				SpecPath: specFooPath,
				Properties: map[string]*SWGSchemaProperty{
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_foo"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_foo":           struct{}{},
							specFooPath + "#/definitions/def_array_ref":     struct{}{},
							specFooPath + "#/definitions/def_array_ref_ref": struct{}{},
						},
					},
				},
				swagger: specFoo,
			},
		},
	}

	for idx, c := range cases {
		actual, err := NewSWGSchema(c.specPath, c.schemaName)
		require.Equal(t, c.err, err, idx)
		require.Equal(t, c.expect, *actual, idx)
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
				*propertyaddr.NewPropertyAddrFromStringWithOwner("def_regular", "prop_array_of_object"),
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
					"prop_array_of_object.prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_foo"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_regular": struct{}{},
							specFooPath + "#/definitions/def_foo":     struct{}{},
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
		require.Equal(t, c.err, err, idx)
		if err == nil {
			for _, addr := range c.expandAddrs {
				require.NoError(t, swgschema.ExpandPropertyOneLevelDeep(addr), idx)
			}
			require.Equal(t, c.expect, *swgschema, idx)
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
					tfPropAddr:  *propertyaddr.NewPropertyAddrFromString("res1:p1"),
					err:         nil,
					expect: SWGSchema{
						Name:     "def_a",
						SpecPath: specFooPath,
						Properties: map[string]*SWGSchemaProperty{
							"prop_primitive": {
								TFLinks: []TFLink{
									{*propertyaddr.NewPropertyAddrFromString("res1:p1")},
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
							"p2": {
								TFLinks: []TFLink{},
								schema:  specFoo.Definitions["def_a"].Properties["p2"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_a": struct{}{},
								},
							},
							"p3": {
								TFLinks: []TFLink{},
								schema:  specFoo.Definitions["def_a"].Properties["p3"],
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
					tfPropAddr:  *propertyaddr.NewPropertyAddrFromString("res2:p1"),
					err:         nil,
					expect: SWGSchema{
						Name:     "def_a",
						SpecPath: specFooPath,
						Properties: map[string]*SWGSchemaProperty{
							"prop_primitive": {
								TFLinks: []TFLink{
									{*propertyaddr.NewPropertyAddrFromString("res1:p1")},
									{*propertyaddr.NewPropertyAddrFromString("res2:p1")},
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
							"p2": {
								TFLinks: []TFLink{},
								schema:  specFoo.Definitions["def_a"].Properties["p2"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_a": struct{}{},
								},
							},
							"p3": {
								TFLinks: []TFLink{},
								schema:  specFoo.Definitions["def_a"].Properties["p3"],
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
					tfPropAddr:  *propertyaddr.NewPropertyAddrFromString("res1:p2"),
					err:         nil,
					expect: SWGSchema{
						Name:     "def_a",
						SpecPath: specFooPath,
						Properties: map[string]*SWGSchemaProperty{
							"prop_primitive": {
								TFLinks: []TFLink{
									{*propertyaddr.NewPropertyAddrFromString("res1:p1")},
									{*propertyaddr.NewPropertyAddrFromString("res2:p1")},
								},
								schema: specFoo.Definitions["def_foo"].Properties["prop_primitive"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_a":   struct{}{},
									specFooPath + "#/definitions/def_foo": struct{}{},
								},
							},
							"p1.prop_primitive": {
								TFLinks: []TFLink{
									{*propertyaddr.NewPropertyAddrFromString("res1:p2")},
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
							"p2": {
								TFLinks: []TFLink{},
								schema:  specFoo.Definitions["def_a"].Properties["p2"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_a": struct{}{},
								},
							},
							"p3": {
								TFLinks: []TFLink{},
								schema:  specFoo.Definitions["def_a"].Properties["p3"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_a": struct{}{},
								},
							},
						},
						swagger: specFoo,
					},
				},
				{
					swgPropAddr: *propertyaddr.NewPropertyAddrFromString("def_a:p3.prop_primitive"),
					tfPropAddr:  *propertyaddr.NewPropertyAddrFromString("res1:p3"),
					err:         nil,
					expect: SWGSchema{
						Name:     "def_a",
						SpecPath: specFooPath,
						Properties: map[string]*SWGSchemaProperty{
							"prop_primitive": {
								TFLinks: []TFLink{
									{*propertyaddr.NewPropertyAddrFromString("res1:p1")},
									{*propertyaddr.NewPropertyAddrFromString("res2:p1")},
								},
								schema: specFoo.Definitions["def_foo"].Properties["prop_primitive"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_a":   struct{}{},
									specFooPath + "#/definitions/def_foo": struct{}{},
								},
							},
							"p1.prop_primitive": {
								TFLinks: []TFLink{
									{*propertyaddr.NewPropertyAddrFromString("res1:p2")},
								},
								schema: specBar.Definitions["def_bar"].Properties["prop_primitive"],
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
							"p2": {
								TFLinks: []TFLink{},
								schema:  specFoo.Definitions["def_a"].Properties["p2"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_a": struct{}{},
								},
							},
							"p3.prop_primitive": {
								TFLinks: []TFLink{
									{*propertyaddr.NewPropertyAddrFromString("res1:p3")},
								},
								schema: specBar.Definitions["def_bar"].Properties["prop_primitive"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_a":   struct{}{},
									specBarPath + "#/definitions/def_bar": struct{}{},
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
		require.NoError(t, err, idx)
		if err == nil {
			for iidx, s := range c.steps {
				err := schema.AddTFLink(s.swgPropAddr, s.tfPropAddr)
				require.Equal(t, s.err, err, fmt.Sprintf("%d.%d", idx, iidx))
				if err == nil {
					require.Equal(t, s.expect, *schema, fmt.Sprintf("%d.%d", idx, iidx))
				}
			}
		}
	}
}

func TestSWGSchema_Marshal(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	specFooPath := filepath.Join(pwd, "testdata", "swagger", "foo.json")

	type process func(schema *SWGSchema)

	cases := []struct {
		specPath   string
		schemaName string
		process    process
		expect     string
	}{
		{
			specPath:   specFooPath,
			schemaName: "def_regular",
			expect: fmt.Sprintf(`{
    "Name": "def_regular",
    "SpecPath": "%s",
    "Properties": {
        "prop_array_of_primitive": {
            "TFLinks": []
        },
        "prop_array_of_object": {
            "TFLinks": []
        },
        "prop_object": {
            "TFLinks": []
        },
        "prop_primitive": {
            "TFLinks": []
        }
    }
}`, specFooPath),
		},
		{
			specPath:   specFooPath,
			schemaName: "def_propInFileRef",
			expect: fmt.Sprintf(`{
    "Name": "def_propInFileRef",
    "SpecPath": "%s",
    "Properties": {
        "prop_inFileRef": {
            "TFLinks": []
        }
    }
}`, specFooPath),
		},
		{
			specPath:   specFooPath,
			schemaName: "def_allOf",
			expect: fmt.Sprintf(`{
    "Name": "def_allOf",
    "SpecPath": "%s",
    "Properties": {
		"prop_nested1": {
            "TFLinks": []
 		},
		"prop_nested2": {
            "TFLinks": []
 		},
		"prop_primitive": {
            "TFLinks": []
 		},
		"p1": {
            "TFLinks": []
 		}
    }
}`, specFooPath),
		},
		{
			specPath:   specFooPath,
			schemaName: "def_a",
			process: func(schema *SWGSchema) {
				swgPropAddr := *propertyaddr.NewPropertyAddrFromString("def_a:p1.prop_primitive")
				tfPropAddr := *propertyaddr.NewPropertyAddrFromString("res1:p2")
				require.NoError(t, schema.AddTFLink(swgPropAddr, tfPropAddr))
			},
			expect: fmt.Sprintf(`{
    "Name": "def_a",
    "SpecPath": "%s",
    "Properties": {
		"prop_primitive": {
            "TFLinks": []
 		},
		"p1.prop_primitive": {
            "TFLinks": ["res1:p2"]
 		},
		"p1.p1_1": {
            "TFLinks": []
 		},
		"p2": {
            "TFLinks": []
 		},
		"p3": {
            "TFLinks": []
 		}
    }
}`, specFooPath),
		},
	}

	for idx, c := range cases {
		schema, err := NewSWGSchema(c.specPath, c.schemaName)
		require.NoError(t, err, idx)
		if c.process != nil {
			c.process(schema)
		}
		b, err := json.Marshal(schema)
		require.NoError(t, err, idx)
		require.JSONEq(t, c.expect, string(b))
	}
}

func TestSWGSchema_Unmarshal(t *testing.T) {
	cases := []struct {
		input  string
		expect SWGSchema
	}{
		{
			input: `{
    "Name": "def_a",
    "SpecPath": "path",
    "Properties": {
		"p1.prop_primitive": {
            "TFLinks": ["res1:p2"]
 		},
		"p1.p1_1": {
            "TFLinks": []
 		},
		"prop_primitive": {
            "TFLinks": []
 		}
    }
}`,
			expect: SWGSchema{
				Name:     "def_a",
				SpecPath: "path",
				Properties: SWGSchemaProperties{
					"p1.prop_primitive": &SWGSchemaProperty{TFLinks: []TFLink{
						{*propertyaddr.NewPropertyAddrFromString("res1:p2")},
					}},
					"p1.p1_1":        &SWGSchemaProperty{TFLinks: []TFLink{}},
					"prop_primitive": &SWGSchemaProperty{TFLinks: []TFLink{}},
				},
			},
		},
	}

	for idx, c := range cases {
		var actual SWGSchema
		require.NoError(t, json.Unmarshal([]byte(c.input), &actual), idx)
		require.Equal(t, c.expect, actual, idx)
	}
}
