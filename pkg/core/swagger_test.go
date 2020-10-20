package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/magodo/terraform-provider-azurerm-insight/pkg/core/propertyaddr"
)

func TestNewSWGSchema(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	specBasePathLocal := filepath.Join(pwd, "testdata", "swagger")
	specFooPathLocal := filepath.Join(specBasePathLocal, "foo.json")
	specBarPathLocal := filepath.Join(specBasePathLocal, "bar.json")
	specFoo, err := LoadSwagger(specFooPathLocal)
	require.NoError(t, err)
	specBar, err := LoadSwagger(specBarPathLocal)
	_ = specBar
	require.NoError(t, err)

	specBaseURL := "https://gist.githubusercontent.com/magodo/f054bb1c2e7a1c74fd78f65eb42a17bb/raw/c1e7508ce27c985390353d7d5f32655028536a13"
	specFooURL := specBaseURL + "/foo.json"
	_ = specFooURL

	cases := []struct {
		specBaseURL string
		specRelPath string
		schemaName  string
		err         error
		expect      SWGSchema
	}{
		// NO.0
		{
			specBaseURL: specBasePathLocal,
			specRelPath: "foo.json",
			schemaName:  "def_foo",
			err:         nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_foo",
				Properties: map[string]*SWGSchemaProperty{
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_foo"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_foo": struct{}{},
						},
					},
				},
				swaggerURL: specFooPathLocal,
				swagger:    specFoo,
			},
		},
		// NO.1
		{
			specBaseURL: specBasePathLocal,
			specRelPath: "foo.json",
			schemaName:  "def_regular",
			err:         nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_regular",
				Properties: map[string]*SWGSchemaProperty{
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_regular"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_regular": struct{}{},
						},
					},
					"prop_array_of_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_regular"].Properties["prop_array_of_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_regular": struct{}{},
						},
					},
					"prop_array_of_object": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_regular"].Properties["prop_array_of_object"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_regular": struct{}{},
						},
					},
					"prop_object": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_regular"].Properties["prop_object"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_regular": struct{}{},
						},
					},
				},
				swaggerURL: specFooPathLocal,
				swagger:    specFoo,
			},
		},
		// NO.2
		{
			specBaseURL: specBasePathLocal,
			specRelPath: "foo.json",
			schemaName:  "def_propInFileRef",
			err:         nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_propInFileRef",
				Properties: map[string]*SWGSchemaProperty{
					"prop_inFileRef": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_propInFileRef"].Properties["prop_inFileRef"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_propInFileRef": struct{}{},
						},
					},
				},
				swaggerURL: specFooPathLocal,
				swagger:    specFoo,
			},
		},
		// NO.3
		{
			specBaseURL: specBasePathLocal,
			specRelPath: "foo.json",
			schemaName:  "def_propSelfRef",
			err:         nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_propSelfRef",
				Properties: map[string]*SWGSchemaProperty{
					"prop_selfRef": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_propSelfRef"].Properties["prop_selfRef"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_propSelfRef": struct{}{},
						},
					},
				},
				swaggerURL: specFooPathLocal,
				swagger:    specFoo,
			},
		},
		// NO.4
		{
			specBaseURL: specBasePathLocal,
			specRelPath: "foo.json",
			schemaName:  "def_propCrossFileRef",
			err:         nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_propCrossFileRef",
				Properties: map[string]*SWGSchemaProperty{
					"prop_crossFileRef": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_propCrossFileRef"].Properties["prop_crossFileRef"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_propCrossFileRef": struct{}{},
						},
					},
				},
				swaggerURL: specFooPathLocal,
				swagger:    specFoo,
			},
		},
		// NO.5
		{
			specBaseURL: specBasePathLocal,
			specRelPath: "foo.json",
			schemaName:  "def_inFileRef",
			err:         nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_inFileRef",
				Properties: map[string]*SWGSchemaProperty{
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_inFileRef"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_inFileRef": struct{}{},
							specFooPathLocal + "#/definitions/def_foo":       struct{}{},
						},
					},
				},
				swaggerURL: specFooPathLocal,
				swagger:    specFoo,
			},
		},
		// NO.6
		{
			specBaseURL: specBasePathLocal,
			specRelPath: "foo.json",
			schemaName:  "def_crossFileRef",
			err:         nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_crossFileRef",
				Properties: map[string]*SWGSchemaProperty{
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specBar.Definitions["def_bar"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_crossFileRef": struct{}{},
							specBarPathLocal + "#/definitions/def_bar":          struct{}{},
						},
					},
				},
				swaggerURL: specFooPathLocal,
				swagger:    specFoo,
			},
		},
		// NO.7
		{
			specBaseURL: specBasePathLocal,
			specRelPath: "foo.json",
			schemaName:  "def_selfRef",
			err:         nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_selfRef",
				Properties: map[string]*SWGSchemaProperty{
					"": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_selfRef"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_selfRef": struct{}{},
						},
					},
				},
				swaggerURL: specFooPathLocal,
				swagger:    specFoo,
			},
		},
		// NO.8
		{
			specBaseURL: specBasePathLocal,
			specRelPath: "foo.json",
			schemaName:  "def_allOf",
			err:         nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_allOf",
				Properties: map[string]*SWGSchemaProperty{
					"p1": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_allOf"].Properties["p1"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_allOf": struct{}{},
						},
					},
					"prop_nested1": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_allOf"].AllOf[0].Properties["prop_nested1"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_allOf": struct{}{},
						},
					},
					"prop_nested2": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_allOf"].AllOf[0].Properties["prop_nested2"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_allOf": struct{}{},
						},
					},
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specBar.Definitions["def_bar"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_allOf": struct{}{},
							specBarPathLocal + "#/definitions/def_bar":   struct{}{},
						},
					},
				},
				swaggerURL: specFooPathLocal,
				swagger:    specFoo,
			},
		},
		// NO.9
		{
			specBaseURL: specBasePathLocal,
			specRelPath: "foo.json",
			schemaName:  "def_array_simple",
			err:         nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_array_simple",
				Properties: map[string]*SWGSchemaProperty{
					"": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_array_simple"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_array_simple": struct{}{},
						},
					},
				},
				swaggerURL: specFooPathLocal,
				swagger:    specFoo,
			},
		},
		// NO.10
		{
			specBaseURL: specBasePathLocal,
			specRelPath: "foo.json",
			schemaName:  "def_array_ref",
			err:         nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_array_ref",
				Properties: map[string]*SWGSchemaProperty{
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_foo"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_foo":       struct{}{},
							specFooPathLocal + "#/definitions/def_array_ref": struct{}{},
						},
					},
				},
				swaggerURL: specFooPathLocal,
				swagger:    specFoo,
			},
		},
		// NO.11
		{
			specBaseURL: specBasePathLocal,
			specRelPath: "foo.json",
			schemaName:  "def_array_ref_ref",
			err:         nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_array_ref_ref",
				Properties: map[string]*SWGSchemaProperty{
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_foo"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_foo":           struct{}{},
							specFooPathLocal + "#/definitions/def_array_ref":     struct{}{},
							specFooPathLocal + "#/definitions/def_array_ref_ref": struct{}{},
						},
					},
				},
				swaggerURL: specFooPathLocal,
				swagger:    specFoo,
			},
		},
		// NO.12: from http swagger spec
		{
			specBaseURL: specBaseURL,
			specRelPath: "foo.json",
			schemaName:  "def_foo",
			err:         nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_foo",
				Properties: map[string]*SWGSchemaProperty{
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_foo"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooURL + "#/definitions/def_foo": struct{}{},
						},
					},
				},
				swaggerURL: specFooURL,
				swagger:    specFoo,
			},
		},
		// NO.13: discriminator
		{
			specBaseURL: specBasePathLocal,
			specRelPath: "foo.json",
			schemaName:  "def_base",
			err:         nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_base",
				Properties: map[string]*SWGSchemaProperty{
					"[def_variant1]": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_variant1"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_variant1": struct{}{},
						},
					},
					"[def_variant2]": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_variant2"],
						resolvedRefs: map[string]interface{}{
							specFooPathLocal + "#/definitions/def_variant2": struct{}{},
						},
					},
				},
				swaggerURL: specFooPathLocal,
				swagger:    specFoo,
			},
		},
	}

	for idx, c := range cases {
		actual, err := NewSWGSchema(c.specBaseURL, c.specRelPath, c.schemaName)
		require.Equal(t, c.err, err, idx)
		require.Equal(t, c.expect, *actual, idx)
	}
}

func TestSWGSchema_ExpandPropertyOneLevelDeep(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	specBasePath := filepath.Join(pwd, "testdata", "swagger")
	specFooPath := filepath.Join(specBasePath, "foo.json")
	specBarPath := filepath.Join(specBasePath, "bar.json")
	specFoo, err := LoadSwagger(specFooPath)
	specBar, err := LoadSwagger(specBarPath)
	_ = specBar
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		swaggerRelPath string
		schemaName     string
		expandAddrs    []propertyaddr.SwaggerPropertyAddr
		err            error
		expect         SWGSchema
	}{
		// NO.0
		{
			swaggerRelPath: "foo.json",
			schemaName:     "def_foo",
			expandAddrs: []propertyaddr.SwaggerPropertyAddr{
				propertyaddr.MustNewSwaggerPropertyAddr("def_foo", "prop_primitive"),
			},
			err: nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_foo",
				Properties: map[string]*SWGSchemaProperty{
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_foo"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_foo": struct{}{},
						},
					},
				},
				swaggerURL: specFooPath,
				swagger:    specFoo,
			},
		},
		// NO.1
		{
			swaggerRelPath: "foo.json",
			schemaName:     "def_foo",
			expandAddrs: []propertyaddr.SwaggerPropertyAddr{
				propertyaddr.MustNewSwaggerPropertyAddr("def_foo", "prop_primitive"),
				propertyaddr.MustNewSwaggerPropertyAddr("def_foo", "prop_primitive"),
			},
			err: nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_foo",
				Properties: map[string]*SWGSchemaProperty{
					"prop_primitive": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_foo"].Properties["prop_primitive"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_foo": struct{}{},
						},
					},
				},
				swaggerURL: specFooPath,
				swagger:    specFoo,
			},
		},
		// NO.2
		{
			swaggerRelPath: "foo.json",
			schemaName:     "def_regular",
			expandAddrs: []propertyaddr.SwaggerPropertyAddr{
				propertyaddr.MustNewSwaggerPropertyAddr("def_regular", "prop_object"),
				propertyaddr.MustNewSwaggerPropertyAddr("def_regular", "prop_array_of_object"),
			},
			err: nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_regular",
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
				swaggerURL: specFooPath,
				swagger:    specFoo,
			},
		},
		// NO.3
		{
			swaggerRelPath: "foo.json",
			schemaName:     "def_propInFileRef",
			expandAddrs: []propertyaddr.SwaggerPropertyAddr{
				propertyaddr.MustNewSwaggerPropertyAddr("def_propInFileRef", "prop_inFileRef"),
			},
			err: nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_propInFileRef",
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
				swaggerURL: specFooPath,
				swagger:    specFoo,
			},
		},
		// NO.4
		{
			swaggerRelPath: "foo.json",
			schemaName:     "def_propSelfRef",
			expandAddrs: []propertyaddr.SwaggerPropertyAddr{
				propertyaddr.MustNewSwaggerPropertyAddr("def_propSelfRef", "prop_selfRef"),
			},
			err: nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_propSelfRef",
				Properties: map[string]*SWGSchemaProperty{
					"prop_selfRef": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_propSelfRef"].Properties["prop_selfRef"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_propSelfRef": struct{}{},
						},
					},
				},
				swaggerURL: specFooPath,
				swagger:    specFoo,
			},
		},
		// NO.5
		{
			swaggerRelPath: "foo.json",
			schemaName:     "def_propCrossFileRef",
			expandAddrs: []propertyaddr.SwaggerPropertyAddr{
				propertyaddr.MustNewSwaggerPropertyAddr("def_propCrossFileRef", "prop_crossFileRef"),
			},
			err: nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_propCrossFileRef",
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
				swaggerURL: specFooPath,
				swagger:    specFoo,
			},
		},
		// NO.6
		{
			swaggerRelPath: "foo.json",
			schemaName:     "def_inFileRef",
			expandAddrs: []propertyaddr.SwaggerPropertyAddr{
				propertyaddr.MustNewSwaggerPropertyAddr("def_inFileRef", "prop_primitive"),
			},
			err: nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_inFileRef",
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
				swaggerURL: specFooPath,
				swagger:    specFoo,
			},
		},
		// NO.7
		{
			swaggerRelPath: "foo.json",
			schemaName:     "def_crossFileRef",
			expandAddrs: []propertyaddr.SwaggerPropertyAddr{
				propertyaddr.MustNewSwaggerPropertyAddr("def_crossFileRef", "prop_primitive"),
			},
			err: nil,
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_crossFileRef",
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
				swaggerURL: specFooPath,
				swagger:    specFoo,
			},
		},
		// NO.8
		{
			swaggerRelPath: "foo.json",
			schemaName:     "def_selfRef",
			err:            nil,
			expandAddrs: []propertyaddr.SwaggerPropertyAddr{
				propertyaddr.MustNewSwaggerPropertyAddr("def_selfRef", ""),
			},
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_selfRef",
				Properties: map[string]*SWGSchemaProperty{
					"": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_selfRef"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_selfRef": struct{}{},
						},
					},
				},
				swaggerURL: specFooPath,
				swagger:    specFoo,
			},
		},
		// NO.9
		{
			swaggerRelPath: "foo.json",
			schemaName:     "def_c",
			err:            nil,
			expandAddrs: []propertyaddr.SwaggerPropertyAddr{
				propertyaddr.MustNewSwaggerPropertyAddr("def_c", "p1"),
			},
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_c",
				Properties: map[string]*SWGSchemaProperty{
					"p1[def_variant1]": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_variant1"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_c":        struct{}{},
							specFooPath + "#/definitions/def_variant1": struct{}{},
						},
					},
					"p1[def_variant2]": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_variant2"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_c":        struct{}{},
							specFooPath + "#/definitions/def_variant2": struct{}{},
						},
					},
				},
				swaggerURL: specFooPath,
				swagger:    specFoo,
			},
		},
		// NO.10
		{
			swaggerRelPath: "foo.json",
			schemaName:     "def_c",
			err:            nil,
			expandAddrs: []propertyaddr.SwaggerPropertyAddr{
				propertyaddr.MustNewSwaggerPropertyAddr("def_c", "p1"),
				propertyaddr.MustNewSwaggerPropertyAddr("def_c", "p1[def_variant1]"),
			},
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "def_c",
				Properties: map[string]*SWGSchemaProperty{
					"p1[def_variant1].type": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_base"].Properties["type"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_c":        struct{}{},
							specFooPath + "#/definitions/def_variant1": struct{}{},
						},
					},
					"p1[def_variant2]": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["def_variant2"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/def_c":        struct{}{},
							specFooPath + "#/definitions/def_variant2": struct{}{},
						},
					},
				},
				swaggerURL: specFooPath,
				swagger:    specFoo,
			},
		},
		// NO.11
		{
			swaggerRelPath: "foo.json",
			schemaName:     "ruleCollectionGroup",
			err:            nil,
			expandAddrs: []propertyaddr.SwaggerPropertyAddr{
				propertyaddr.MustNewSwaggerPropertyAddr("ruleCollectionGroup", "ruleCollections"),
			},
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "ruleCollectionGroup",
				Properties: map[string]*SWGSchemaProperty{
					"ruleCollections[natRuleCollection]": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["natRuleCollection"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/ruleCollectionGroup": struct{}{},
							specFooPath + "#/definitions/natRuleCollection":   struct{}{},
						},
					},
					"ruleCollections[filterRuleCollection]": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["filterRuleCollection"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/ruleCollectionGroup":  struct{}{},
							specFooPath + "#/definitions/filterRuleCollection": struct{}{},
						},
					},
				},
				swaggerURL: specFooPath,
				swagger:    specFoo,
			},
		},
		// NO.12
		{
			swaggerRelPath: "foo.json",
			schemaName:     "ruleCollectionGroup",
			err:            nil,
			expandAddrs: []propertyaddr.SwaggerPropertyAddr{
				propertyaddr.MustNewSwaggerPropertyAddr("ruleCollectionGroup", "ruleCollections"),
				propertyaddr.MustNewSwaggerPropertyAddr("ruleCollectionGroup", "ruleCollections[natRuleCollection]"),
			},
			expect: SWGSchema{
				SwaggerRelPath: "foo.json",
				Name:           "ruleCollectionGroup",
				Properties: map[string]*SWGSchemaProperty{
					"ruleCollections[natRuleCollection].action": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["natRuleCollection"].Properties["action"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/ruleCollectionGroup": struct{}{},
							specFooPath + "#/definitions/natRuleCollection":   struct{}{},
						},
					},
					"ruleCollections[natRuleCollection].ruleCollectionType": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["ruleCollection"].Properties["ruleCollectionType"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/ruleCollectionGroup": struct{}{},
							specFooPath + "#/definitions/natRuleCollection":   struct{}{},
						},
					},
					"ruleCollections[natRuleCollection].name": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["ruleCollection"].Properties["name"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/ruleCollectionGroup": struct{}{},
							specFooPath + "#/definitions/natRuleCollection":   struct{}{},
						},
					},
					"ruleCollections[filterRuleCollection]": {
						TFLinks: []TFLink{},
						schema:  specFoo.Definitions["filterRuleCollection"],
						resolvedRefs: map[string]interface{}{
							specFooPath + "#/definitions/ruleCollectionGroup":  struct{}{},
							specFooPath + "#/definitions/filterRuleCollection": struct{}{},
						},
					},
				},
				swaggerURL: specFooPath,
				swagger:    specFoo,
			},
		},
	}

	for idx, c := range cases {
		swgschema, err := NewSWGSchema(specBasePath, c.swaggerRelPath, c.schemaName)
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
	specBasePath := filepath.Join(pwd, "testdata", "swagger")
	specFooPath := filepath.Join(specBasePath, "foo.json")
	specBarPath := filepath.Join(specBasePath, "bar.json")
	specFoo, err := LoadSwagger(specFooPath)
	specBar, err := LoadSwagger(specBarPath)
	_ = specBar
	if err != nil {
		t.Fatal(err)
	}

	type step struct {
		swgPropAddr propertyaddr.SwaggerPropertyAddr
		tfPropAddr  propertyaddr.TerraformPropertyAddr
		err         bool
		expect      SWGSchema
	}

	cases := []struct {
		swaggerRelPath string
		schemaName     string
		steps          []step
	}{
		{
			swaggerRelPath: "foo.json",
			schemaName:     "def_a",
			steps: []step{
				{
					swgPropAddr: propertyaddr.MustParseSwaggerPropertyAddr("def_a:prop_primitive"),
					tfPropAddr:  *propertyaddr.ParseTerraformPropertyAddr("res1:p1"),
					expect: SWGSchema{
						SwaggerRelPath: "foo.json",
						Name:           "def_a",
						Properties: map[string]*SWGSchemaProperty{
							"prop_primitive": {
								TFLinks: []TFLink{
									{*propertyaddr.ParseTerraformPropertyAddr("res1:p1")},
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
						swaggerURL: specFooPath,
						swagger:    specFoo,
					},
				},
				// add a second tf link to the same swg property
				{
					swgPropAddr: propertyaddr.MustParseSwaggerPropertyAddr("def_a:prop_primitive"),
					tfPropAddr:  *propertyaddr.ParseTerraformPropertyAddr("res2:p1"),
					expect: SWGSchema{
						SwaggerRelPath: "foo.json",
						Name:           "def_a",
						Properties: map[string]*SWGSchemaProperty{
							"prop_primitive": {
								TFLinks: []TFLink{
									{*propertyaddr.ParseTerraformPropertyAddr("res1:p1")},
									{*propertyaddr.ParseTerraformPropertyAddr("res2:p1")},
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
						swaggerURL: specFooPath,
						swagger:    specFoo,
					},
				},
				{
					swgPropAddr: propertyaddr.MustParseSwaggerPropertyAddr("def_a:p1.prop_primitive"),
					tfPropAddr:  *propertyaddr.ParseTerraformPropertyAddr("res1:p2"),
					expect: SWGSchema{
						SwaggerRelPath: "foo.json",
						Name:           "def_a",
						Properties: map[string]*SWGSchemaProperty{
							"prop_primitive": {
								TFLinks: []TFLink{
									{*propertyaddr.ParseTerraformPropertyAddr("res1:p1")},
									{*propertyaddr.ParseTerraformPropertyAddr("res2:p1")},
								},
								schema: specFoo.Definitions["def_foo"].Properties["prop_primitive"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_a":   struct{}{},
									specFooPath + "#/definitions/def_foo": struct{}{},
								},
							},
							"p1.prop_primitive": {
								TFLinks: []TFLink{
									{*propertyaddr.ParseTerraformPropertyAddr("res1:p2")},
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
						swaggerURL: specFooPath,
						swagger:    specFoo,
					},
				},
				{
					swgPropAddr: propertyaddr.MustParseSwaggerPropertyAddr("def_a:p3.prop_primitive"),
					tfPropAddr:  *propertyaddr.ParseTerraformPropertyAddr("res1:p3"),
					expect: SWGSchema{
						SwaggerRelPath: "foo.json",
						Name:           "def_a",
						Properties: map[string]*SWGSchemaProperty{
							"prop_primitive": {
								TFLinks: []TFLink{
									{*propertyaddr.ParseTerraformPropertyAddr("res1:p1")},
									{*propertyaddr.ParseTerraformPropertyAddr("res2:p1")},
								},
								schema: specFoo.Definitions["def_foo"].Properties["prop_primitive"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_a":   struct{}{},
									specFooPath + "#/definitions/def_foo": struct{}{},
								},
							},
							"p1.prop_primitive": {
								TFLinks: []TFLink{
									{*propertyaddr.ParseTerraformPropertyAddr("res1:p2")},
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
									{*propertyaddr.ParseTerraformPropertyAddr("res1:p3")},
								},
								schema: specBar.Definitions["def_bar"].Properties["prop_primitive"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_a":   struct{}{},
									specBarPath + "#/definitions/def_bar": struct{}{},
								},
							},
						},
						swaggerURL: specFooPath,
						swagger:    specFoo,
					},
				},
			},
		},

		// add tf link to swagger child property first, then add another tf link to swagger parent property
		{
			swaggerRelPath: "foo.json",
			schemaName:     "def_b",
			steps: []step{
				{
					swgPropAddr: propertyaddr.MustParseSwaggerPropertyAddr("def_b:p1.p1_1"),
					tfPropAddr:  *propertyaddr.ParseTerraformPropertyAddr("res2:p1"),
					expect: SWGSchema{
						SwaggerRelPath: "foo.json",
						Name:           "def_b",
						Properties: map[string]*SWGSchemaProperty{
							"p1.p1_1": {
								TFLinks: []TFLink{
									{*propertyaddr.ParseTerraformPropertyAddr("res2:p1")},
								},
								schema: specFoo.Definitions["def_b"].Properties["p1"].Properties["p1_1"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_b": struct{}{},
								},
							},
						},
						swaggerURL: specFooPath,
						swagger:    specFoo,
					},
				},
				{
					swgPropAddr: propertyaddr.MustParseSwaggerPropertyAddr("def_b:p1"),
					tfPropAddr:  *propertyaddr.ParseTerraformPropertyAddr("res1:p1"),
					expect: SWGSchema{
						SwaggerRelPath: "foo.json",
						Name:           "def_b",
						Properties: map[string]*SWGSchemaProperty{
							"p1.p1_1": {
								TFLinks: []TFLink{
									{*propertyaddr.ParseTerraformPropertyAddr("res2:p1")},
									{*propertyaddr.ParseTerraformPropertyAddr("res1:p1")},
								},
								schema: specFoo.Definitions["def_b"].Properties["p1"].Properties["p1_1"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_b": struct{}{},
								},
							},
						},
						swaggerURL: specFooPath,
						swagger:    specFoo,
					},
				},
			},
		},

		// add tf link to swagger parent property first, then add another tf link to swagger child property
		{
			swaggerRelPath: "foo.json",
			schemaName:     "def_b",
			steps: []step{
				{
					swgPropAddr: propertyaddr.MustParseSwaggerPropertyAddr("def_b:p1"),
					tfPropAddr:  *propertyaddr.ParseTerraformPropertyAddr("res1:p1"),
					expect: SWGSchema{
						SwaggerRelPath: "foo.json",
						Name:           "def_b",
						Properties: map[string]*SWGSchemaProperty{
							"p1": {
								TFLinks: []TFLink{
									{*propertyaddr.ParseTerraformPropertyAddr("res1:p1")},
								},
								schema: specFoo.Definitions["def_b"].Properties["p1"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_b": struct{}{},
								},
							},
						},
						swaggerURL: specFooPath,
						swagger:    specFoo,
					},
				},
				{
					swgPropAddr: propertyaddr.MustParseSwaggerPropertyAddr("def_b:p1.p1_1"),
					tfPropAddr:  *propertyaddr.ParseTerraformPropertyAddr("res2:p1"),
					expect: SWGSchema{
						SwaggerRelPath: "foo.json",
						Name:           "def_b",
						Properties: map[string]*SWGSchemaProperty{
							"p1.p1_1": {
								TFLinks: []TFLink{
									{*propertyaddr.ParseTerraformPropertyAddr("res1:p1")},
									{*propertyaddr.ParseTerraformPropertyAddr("res2:p1")},
								},
								schema: specFoo.Definitions["def_b"].Properties["p1"].Properties["p1_1"],
								resolvedRefs: map[string]interface{}{
									specFooPath + "#/definitions/def_b": struct{}{},
								},
							},
						},
						swaggerURL: specFooPath,
						swagger:    specFoo,
					},
				},
			},
		},
	}

	for idx, c := range cases {
		schema, err := NewSWGSchema(specBasePath, c.swaggerRelPath, c.schemaName)
		require.NoError(t, err, idx)
		if err == nil {
			for iidx, s := range c.steps {
				if s.err {
					require.Error(t, schema.AddTFLink(s.swgPropAddr, s.tfPropAddr), fmt.Sprintf("%d.%d", idx, iidx))
					continue
				}
				require.NoError(t, schema.AddTFLink(s.swgPropAddr, s.tfPropAddr), fmt.Sprintf("%d.%d", idx, iidx))
				require.Equal(t, s.expect, *schema, fmt.Sprintf("%d.%d", idx, iidx))
			}
		}
	}
}

func TestSWGSchema_Marshal(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	specBasePath := filepath.Join(pwd, "testdata", "swagger")

	type process func(schema *SWGSchema)

	cases := []struct {
		swaggerRelPath string
		schemaName     string
		process        process
		expect         string
	}{
		{
			swaggerRelPath: "foo.json",
			schemaName:     "def_regular",
			expect: fmt.Sprintf(`{
    "SwaggerRelPath": "foo.json",
    "Name": "def_regular",
    "Properties": {
        "prop_array_of_primitive": {},
        "prop_array_of_object": {},
        "prop_object": {},
        "prop_primitive": {}
    }
}`),
		},
		{
			swaggerRelPath: "foo.json",
			schemaName:     "def_propInFileRef",
			expect: fmt.Sprintf(`{
    "SwaggerRelPath": "foo.json",
    "Name": "def_propInFileRef",
    "Properties": {
        "prop_inFileRef": {}
    }
}`),
		},
		{
			swaggerRelPath: "foo.json",
			schemaName:     "def_allOf",
			expect: fmt.Sprintf(`{
    "SwaggerRelPath": "foo.json",
    "Name": "def_allOf",
    "Properties": {
		"prop_nested1": {},
		"prop_nested2": {},
		"prop_primitive": {},
		"p1": {}
    }
}`),
		},
		{
			swaggerRelPath: "foo.json",
			schemaName:     "def_a",
			process: func(schema *SWGSchema) {
				swgPropAddr := propertyaddr.MustParseSwaggerPropertyAddr("def_a:p1.prop_primitive")
				tfPropAddr := *propertyaddr.ParseTerraformPropertyAddr("res1:p2")
				require.NoError(t, schema.AddTFLink(swgPropAddr, tfPropAddr))
			},
			expect: fmt.Sprintf(`{
    "SwaggerRelPath": "foo.json",
    "Name": "def_a",
    "Properties": {
		"prop_primitive": {},
		"p1.prop_primitive": {
            "TFLinks": ["res1:p2"]
 		},
		"p1.p1_1": {},
		"p2": {},
		"p3": {}
    }
}`),
		},
	}

	for idx, c := range cases {
		schema, err := NewSWGSchema(specBasePath, c.swaggerRelPath, c.schemaName)
		require.NoError(t, err, idx)
		if c.process != nil {
			c.process(schema)
		}
		b, err := json.Marshal(schema)
		require.NoError(t, err, idx)
		require.JSONEq(t, c.expect, string(b), idx)
	}
}

func TestSWGSchema_Unmarshal(t *testing.T) {
	cases := []struct {
		input  string
		expect SWGSchema
	}{
		{
			input: `{
    "SwaggerRelPath": "foo.json",
    "Name": "def_a",
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
				SwaggerRelPath: "foo.json",
				Name:           "def_a",
				Properties: SWGSchemaProperties{
					"p1.prop_primitive": &SWGSchemaProperty{TFLinks: []TFLink{
						{*propertyaddr.ParseTerraformPropertyAddr("res1:p2")},
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

func TestSWGSchema_CalcCoverage(t *testing.T) {
	cases := []struct {
		swgschema   SWGSchema
		expectStore SWGPropertyCoverageStore
	}{
		{
			swgschema: SWGSchema{
				Properties: SWGSchemaProperties{
					"prop1.covered": {
						TFLinks: []TFLink{{}},
					},
					"prop1.not_covered": {
						TFLinks: []TFLink{},
					},
					"prop2.covered": {
						TFLinks: []TFLink{{}},
					},
					"prop_granted": {
						IsGranted: true,
						TFLinks:   []TFLink{},
					},
				},
			},
			expectStore: SWGPropertyCoverageStore{
				node: swgPropertyCoverageNode{
					TotalAmount:   3,
					CoveredAmount: 2,
					Children: map[string]*swgPropertyCoverageNode{
						"prop1": {
							TotalAmount:   2,
							CoveredAmount: 1,
							Children: map[string]*swgPropertyCoverageNode{
								"not_covered": {
									TotalAmount:   1,
									CoveredAmount: 0,
									Children:      map[string]*swgPropertyCoverageNode{},
								},
								"covered": {
									TotalAmount:   1,
									CoveredAmount: 1,
									Children:      map[string]*swgPropertyCoverageNode{},
								},
							},
						},
						"prop2": {
							TotalAmount:   1,
							CoveredAmount: 1,
							Children: map[string]*swgPropertyCoverageNode{
								"covered": {
									TotalAmount:   1,
									CoveredAmount: 1,
									Children:      map[string]*swgPropertyCoverageNode{},
								},
							},
						},
					},
				},
			},
		},
	}

	for idx, c := range cases {
		require.NoError(t, c.swgschema.CalcCoverage(), idx)
		require.Equal(t, c.expectStore, c.swgschema.coverageStore, idx)
	}
}

func TestSWGSchemas_Grant(t *testing.T) {
	cases := []struct {
		swggrant         SWGGrant
		swgschemas       SWGSchemas
		expectSwgSchemas SWGSchemas
		expectError      bool
	}{
		// grant schema
		{
			swggrant: map[SWGSchemaAddr]SWGSchemaGrant{
				NewSWGSchemaAddr("swaggerRelPath", "schema1"): {
					Comment: "granted because of some reason",
				},
			},
			swgschemas: SWGSchemas{
				m: map[SWGSchemaAddr]*SWGSchema{
					NewSWGSchemaAddr("swaggerRelPath", "schema1"): {
						SwaggerRelPath: "swaggerRelPath",
						Name:           "schema1",
					},
				},
			},
			expectSwgSchemas: SWGSchemas{
				m: map[SWGSchemaAddr]*SWGSchema{
					NewSWGSchemaAddr("swaggerRelPath", "schema1"): {
						IsGranted:      true,
						GrantComment:   "granted because of some reason",
						SwaggerRelPath: "swaggerRelPath",
						Name:           "schema1",
					},
				},
			},
		},
		// grant property
		{
			swggrant: map[SWGSchemaAddr]SWGSchemaGrant{
				NewSWGSchemaAddr("swaggerRelPath", "schema1"): {
					Properties: map[string]string{
						"prop1": "granted because of some reason",
					},
				},
			},
			swgschemas: SWGSchemas{
				m: map[SWGSchemaAddr]*SWGSchema{
					NewSWGSchemaAddr("swaggerRelPath", "schema1"): {
						SwaggerRelPath: "swaggerRelPath",
						Name:           "schema1",
						Properties: map[string]*SWGSchemaProperty{
							"prop1": {
								TFLinks: []TFLink{},
							},
							"prop2": {
								TFLinks: []TFLink{},
							},
						},
					},
				},
			},
			expectSwgSchemas: SWGSchemas{
				m: map[SWGSchemaAddr]*SWGSchema{
					NewSWGSchemaAddr("swaggerRelPath", "schema1"): {
						SwaggerRelPath: "swaggerRelPath",
						Name:           "schema1",
						Properties: map[string]*SWGSchemaProperty{
							"prop1": {
								IsGranted:    true,
								GrantComment: "granted because of some reason",
								TFLinks:      []TFLink{},
							},
							"prop2": {
								TFLinks: []TFLink{},
							},
						},
					},
				},
			},
		},

		// the property to be granted doesn't exist
		{
			swggrant: map[SWGSchemaAddr]SWGSchemaGrant{
				NewSWGSchemaAddr("swaggerRelPath", "schema1"): {
					Properties: map[string]string{
						"non_exist_prop1": "granted because of some reason",
					},
				},
			},
			swgschemas: SWGSchemas{
				m: map[SWGSchemaAddr]*SWGSchema{
					NewSWGSchemaAddr("swaggerRelPath", "schema1"): {
						SwaggerRelPath: "swaggerRelPath",
						Name:           "schema1",
						Properties: map[string]*SWGSchemaProperty{
							"prop1": {
								TFLinks: []TFLink{},
							},
						},
					},
				},
			},
			expectError: true,
		},
	}
	for idx, c := range cases {
		swgschemas := c.swgschemas
		if !c.expectError {
			require.NoError(t, swgschemas.Grant(c.swggrant), idx)
			require.Equal(t, c.expectSwgSchemas, swgschemas, idx)
		} else {
			require.Error(t, swgschemas.Grant(c.swggrant), idx)
		}
	}
}
