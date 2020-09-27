package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/magodo/terraform-provider-azurerm-insight/pkg/core/propertyaddr"
	"github.com/stretchr/testify/assert"
)

func TestNewSchemaScaffoldFromTerraformBlock(t *testing.T) {
	jsonInput := []byte(`
{
  "attributes": {
	"foo": {},
	"bar": {
	  "type": [
		"list",
		[
		  "object",
		  {
			"p1": "string",
			"p2": [
              "list",
              [
              	"object",
              	{
                  "p2_1": "string"
              	}
              ]
			],
			"p3": [
			  "map",
			  "string"
			]
		  }
		]
	  ]
	},
	"baz": {
	  "type": [
		"map",
		"string"
	  ]
	}
  },
  "block_types": {
    "block_a": {
      "block": {
        "attributes": {
          "foo": {}
        },
	    "block_types": {
		  "block_a_a": {
		    "block": {
			  "attributes": {
			    "bar": {}
			  }
			}
		  }
	 	}
      }
    }
  }
}`)

	expect := &TFSchema{
		Name: "res1",
		PropertyLinks: map[string][]SwaggerLink{
			"foo":                   {},
			"bar.p1":                {},
			"bar.p2.p2_1":           {},
			"bar.p3":                {},
			"baz":                   {},
			"block_a.block_a_a.bar": {},
			"block_a.foo":           {},
		},
	}

	var block TerraformBlock
	if err := json.Unmarshal(jsonInput, &block); err != nil {
		t.Fatal(err)
	}
	schema := NewSchemaScaffoldFromTerraformBlock("res1", &block)

	assert.Equal(t, *expect, *schema)
}

func TestMarshalTFSchema(t *testing.T) {
	tfschema := TFSchema{
		Name:        "res1",
		SwaggerSpec: "spec1",
		PropertyLinks: map[string][]SwaggerLink{
			"bar": {
				{
					Spec:       strPtr("xxx"),
					SchemaProp: *propertyaddr.NewPropertyAddrFromString("schema1:p1.p2"),
				},
				{
					Spec:       strPtr("yyy"),
					SchemaProp: *propertyaddr.NewPropertyAddrFromString("schema2:p3.p4"),
				},
			},
			"block_a.block_a_a.bar": {
				{
					Spec:       strPtr("xxx"),
					SchemaProp: *propertyaddr.NewPropertyAddrFromString("schema1:p1.p2"),
				},
				{
					Spec:       strPtr("yyy"),
					SchemaProp: *propertyaddr.NewPropertyAddrFromString("schema2:p3.p4"),
				},
			},
		},
	}

	expect := `{
    "Name": "res1",
    "PropertyLinks": {
        "bar": [
            {
                "prop": "schema1:p1.p2",
                "swagger": "xxx"
            },
            {
                "prop": "schema2:p3.p4",
                "swagger": "yyy"
            }
        ],
        "block_a.block_a_a.bar": [
            {
                "prop": "schema1:p1.p2",
                "swagger": "xxx"
            },
            {
                "prop": "schema2:p3.p4",
                "swagger": "yyy"
            }
        ]
    },
    "swagger": "spec1"
}`

	actual, err := json.Marshal(tfschema)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, expect, string(actual))
}

func TestUnmarshalTFSchema(t *testing.T) {
	jsonInput := []byte(`{
    "Name": "res1",
    "PropertyLinks": {
        "bar": [
            {
                "prop": "schema1:p1.p2",
                "swagger": "xxx"
            },
            {
                "prop": "schema2:p3.p4",
                "swagger": "yyy"
            }
        ],
        "block_a::block_a_a::bar": [
            {
                "prop": "schema1:p1.p2",
                "swagger": "xxx"
            },
            {
                "prop": "schema2:p3.p4",
                "swagger": "yyy"
            }
        ]
    },
    "swagger": "spec1"
}`)

	expect := TFSchema{
		Name:        "res1",
		SwaggerSpec: "spec1",
		PropertyLinks: map[string][]SwaggerLink{
			"bar": {
				{
					Spec:       strPtr("xxx"),
					SchemaProp: *propertyaddr.NewPropertyAddrFromString("schema1:p1.p2"),
				},
				{
					Spec:       strPtr("yyy"),
					SchemaProp: *propertyaddr.NewPropertyAddrFromString("schema2:p3.p4"),
				},
			},
			"block_a::block_a_a::bar": {
				{
					Spec:       strPtr("xxx"),
					SchemaProp: *propertyaddr.NewPropertyAddrFromString("schema1:p1.p2"),
				},
				{
					Spec:       strPtr("yyy"),
					SchemaProp: *propertyaddr.NewPropertyAddrFromString("schema2:p3.p4"),
				},
			},
		},
	}

	var schema TFSchema
	if err := json.Unmarshal(jsonInput, &schema); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expect, schema)
}

func TestTFSchema_Validate(t *testing.T) {
	cases := []struct {
		schema TFSchema
		err    error
	}{
		{
			schema: TFSchema{
				Name:        "foo",
				SwaggerSpec: "spec1",
				PropertyLinks: map[string][]SwaggerLink{
					"p1": {
						{
							Spec:       strPtr("spec2"),
							SchemaProp: *propertyaddr.NewPropertyAddrFromString("schema1:p1"),
						},
					},
				},
			},
			err: nil,
		},
		{
			schema: TFSchema{
				Name:        "foo",
				SwaggerSpec: "spec1",
				PropertyLinks: map[string][]SwaggerLink{
					"foo:p1": {
						{
							Spec:       strPtr("spec2"),
							SchemaProp: *propertyaddr.NewPropertyAddrFromString("schema1:p1"),
						},
					},
				},
			},
			err: errors.New("terraform property addr foo:p1 should not specify owner"),
		},
		{
			schema: TFSchema{
				Name:        "foo",
				SwaggerSpec: "spec1",
				PropertyLinks: map[string][]SwaggerLink{
					"p1": {
						{
							Spec:       strPtr("spec2"),
							SchemaProp: *propertyaddr.NewPropertyAddrFromString("p1.p2"),
						},
					},
				},
			},
			err: errors.New("swagger property addr p1.p2 should specify owner"),
		},
	}

	for idx, c := range cases {
		err := c.schema.Validate()
		if c.err != nil {
			assert.EqualError(t, err, c.err.Error(), idx)
		} else {
			assert.NoError(t, err, idx)
		}
	}
}

func TestTFSchema_LinkSwagger(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	specBasePath := filepath.Join(pwd, "testdata", "swagger")

	cases := []struct {
		schemas []TFSchema
		expect  map[string]*SWGSchema
	}{
		// single tf schema -> single swagger
		{
			[]TFSchema{
				{
					Name:        "res1",
					SwaggerSpec: "foo.json",
					PropertyLinks: map[string][]SwaggerLink{
						"p1": {
							{
								SchemaProp: *propertyaddr.NewPropertyAddrFromString("def_a:prop_primitive"),
							},
						},
						"p2.p2_1": {
							{
								SchemaProp: *propertyaddr.NewPropertyAddrFromString("def_a:p1"),
							},
						},
						"p3": {
							{
								SchemaProp: *propertyaddr.NewPropertyAddrFromString("def_a:p2"),
							},
						},
						"p4.p4_1": {
							{
								SchemaProp: *propertyaddr.NewPropertyAddrFromString("def_a:p3.prop_primitive"),
							},
						},
					},
				},
			},
			map[string]*SWGSchema{
				"foo.json" + "#/definitions/def_a": {
					SwaggerRelPath: "foo.json",
					Name:           "def_a",
					Properties: map[string]*SWGSchemaProperty{
						"prop_primitive": {
							TFLinks: []TFLink{
								{
									*propertyaddr.NewPropertyAddrFromString("res1:p1"),
								},
							},
						},
						"p1": {
							TFLinks: []TFLink{
								{
									*propertyaddr.NewPropertyAddrFromString("res1:p2.p2_1"),
								},
							},
						},
						"p2": {TFLinks: []TFLink{
							{
								*propertyaddr.NewPropertyAddrFromString("res1:p3"),
							},
						}},
						"p3.prop_primitive": {TFLinks: []TFLink{
							{
								*propertyaddr.NewPropertyAddrFromString("res1:p4.p4_1"),
							},
						}},
					},
				},
			},
		},
		// single tf schema -> multiple swaggers (cross file)
		{
			[]TFSchema{
				{
					Name:        "res1",
					SwaggerSpec: "foo.json",
					PropertyLinks: map[string][]SwaggerLink{
						"p1": {
							{
								SchemaProp: *propertyaddr.NewPropertyAddrFromString("def_a:prop_primitive"),
							},
						},
						"p2.p2_1": {
							{
								SchemaProp: *propertyaddr.NewPropertyAddrFromString("def_a:p1"),
							},
						},
						"p3": {
							{
								Spec:       strPtr("bar.json"),
								SchemaProp: *propertyaddr.NewPropertyAddrFromString("def_bar:prop_primitive"),
							},
						},
					},
				},
			},
			map[string]*SWGSchema{
				"foo.json" + "#/definitions/def_a": {
					Name:           "def_a",
					SwaggerRelPath: "foo.json",
					Properties: map[string]*SWGSchemaProperty{
						"prop_primitive": {
							TFLinks: []TFLink{
								{
									*propertyaddr.NewPropertyAddrFromString("res1:p1"),
								},
							},
						},
						"p1": {
							TFLinks: []TFLink{
								{
									*propertyaddr.NewPropertyAddrFromString("res1:p2.p2_1"),
								},
							},
						},
						"p2": {TFLinks: []TFLink{}},
						"p3": {TFLinks: []TFLink{}},
					},
				},
				"bar.json" + "#/definitions/def_bar": {
					Name:           "def_bar",
					SwaggerRelPath: "bar.json",
					Properties: map[string]*SWGSchemaProperty{
						"prop_primitive": {
							TFLinks: []TFLink{
								{
									*propertyaddr.NewPropertyAddrFromString("res1:p3"),
								},
							},
						},
					},
				},
			},
		},
		// multiple tf schema -> single swaggers
		{
			[]TFSchema{
				{
					Name:        "res1",
					SwaggerSpec: "foo.json",
					PropertyLinks: map[string][]SwaggerLink{
						"p1": {
							{
								SchemaProp: *propertyaddr.NewPropertyAddrFromString("def_a:prop_primitive"),
							},
						},
						"p2.p2_1": {
							{
								SchemaProp: *propertyaddr.NewPropertyAddrFromString("def_a:p1"),
							},
						},
					},
				},
				{
					Name:        "res2",
					SwaggerSpec: "foo.json",
					PropertyLinks: map[string][]SwaggerLink{
						"p1": {
							{
								SchemaProp: *propertyaddr.NewPropertyAddrFromString("def_a:p1.p1_1"),
							},
						},
					},
				},
			},
			map[string]*SWGSchema{
				"foo.json" + "#/definitions/def_a": {
					Name:           "def_a",
					SwaggerRelPath: "foo.json",
					Properties: map[string]*SWGSchemaProperty{
						"prop_primitive": {
							TFLinks: []TFLink{
								{
									*propertyaddr.NewPropertyAddrFromString("res1:p1"),
								},
							},
						},
						"p1.prop_primitive": {
							TFLinks: []TFLink{
								{
									*propertyaddr.NewPropertyAddrFromString("res1:p2.p2_1"),
								},
							},
						},
						"p1.p1_1": {
							TFLinks: []TFLink{
								{
									*propertyaddr.NewPropertyAddrFromString("res1:p2.p2_1"),
								},
								{
									*propertyaddr.NewPropertyAddrFromString("res2:p1"),
								},
							},
						},
						"p2": {TFLinks: []TFLink{}},
						"p3": {TFLinks: []TFLink{}},
					},
				},
			},
		},
		// multiple tf schema -> multiple swaggers
		{
			[]TFSchema{
				{
					Name:        "res1",
					SwaggerSpec: "foo.json",
					PropertyLinks: map[string][]SwaggerLink{
						"p1": {
							{
								SchemaProp: *propertyaddr.NewPropertyAddrFromString("def_a:prop_primitive"),
							},
						},
						"p2.p2_1": {
							{
								SchemaProp: *propertyaddr.NewPropertyAddrFromString("def_a:p1"),
							},
						},
						"p3": {
							{
								Spec:       strPtr("bar.json"),
								SchemaProp: *propertyaddr.NewPropertyAddrFromString("def_bar:prop_primitive"),
							},
						},
					},
				},
				{
					Name:        "res2",
					SwaggerSpec: "foo.json",
					PropertyLinks: map[string][]SwaggerLink{
						"p1": {
							{
								SchemaProp: *propertyaddr.NewPropertyAddrFromString("def_a:p1.p1_1"),
							},
						},
					},
				},
			},
			map[string]*SWGSchema{
				"foo.json" + "#/definitions/def_a": {
					Name:           "def_a",
					SwaggerRelPath: "foo.json",
					Properties: map[string]*SWGSchemaProperty{
						"prop_primitive": {
							TFLinks: []TFLink{
								{
									*propertyaddr.NewPropertyAddrFromString("res1:p1"),
								},
							},
						},
						"p1.prop_primitive": {
							TFLinks: []TFLink{
								{
									*propertyaddr.NewPropertyAddrFromString("res1:p2.p2_1"),
								},
							},
						},
						"p1.p1_1": {
							TFLinks: []TFLink{
								{
									*propertyaddr.NewPropertyAddrFromString("res1:p2.p2_1"),
								},
								{
									*propertyaddr.NewPropertyAddrFromString("res2:p1"),
								},
							},
						},
						"p2": {TFLinks: []TFLink{}},
						"p3": {TFLinks: []TFLink{}},
					},
				},
				"bar.json" + "#/definitions/def_bar": {
					Name:           "def_bar",
					SwaggerRelPath: "bar.json",
					Properties: map[string]*SWGSchemaProperty{
						"prop_primitive": {
							TFLinks: []TFLink{
								{
									*propertyaddr.NewPropertyAddrFromString("res1:p3"),
								},
							},
						},
					},
				},
			},
		},
	}

	for idx, c := range cases {
		swgschemas := NewSGWSchemas()
		for iidx, schema := range c.schemas {
			require.NoError(t, schema.LinkSwagger(swgschemas, specBasePath), fmt.Sprintf("%d.%d", idx, iidx))
		}
		var actual map[string]*SWGSchema
		b, err := json.Marshal(swgschemas.GetAll())
		require.NoError(t, err, idx)
		require.NoError(t, json.Unmarshal(b, &actual), idx)
		require.Equal(t, c.expect, actual, idx)
	}
}

func strPtr(s string) *string {
	return &s
}
