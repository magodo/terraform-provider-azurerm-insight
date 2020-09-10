package pkg

import (
	"encoding/json"
	"errors"
	"github.com/magodo/terraform-provider-azurerm-insight/pkg/propertyaddr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSchemaScaffoldFromTerraformBlock(t *testing.T) {
	jsonInput := []byte(`
{
  "attributes": {
	"foo": {},
    "bar": {}
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
			"bar":                   []SwaggerLink{},
			"block_a.block_a_a.bar": []SwaggerLink{},
			"block_a.foo":           []SwaggerLink{},
			"foo":                   []SwaggerLink{},
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

	expect := `{
    "Name": "res1",
    "Properties": {
        "bar": [
            {
                "prop": {
                    "addr": "p1.p2",
                    "owner": "schema1"
                },
                "swagger": "xxx"
            },
            {
                "prop": {
                    "addr": "p3.p4",
                    "owner": "schema2"
                },
                "swagger": "yyy"
            }
        ],
        "block_a::block_a_a::bar": [
            {
                "prop": {
                    "addr": "p1.p2",
                    "owner": "schema1"
                },
                "swagger": "xxx"
            },
            {
                "prop": {
                    "addr": "p3.p4",
                    "owner": "schema2"
                },
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
    "Properties": {
        "bar": [
            {
                "prop": {
                    "addr": "p1.p2",
                    "owner": "schema1"
                },
                "swagger": "xxx"
            },
            {
                "prop": {
                    "addr": "p3.p4",
                    "owner": "schema2"
                },
                "swagger": "yyy"
            }
        ],
        "block_a::block_a_a::bar": [
            {
                "prop": {
                    "addr": "p1.p2",
                    "owner": "schema1"
                },
                "swagger": "xxx"
            },
            {
                "prop": {
                    "addr": "p3.p4",
                    "owner": "schema2"
                },
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
	cases := []struct{
		schema TFSchema
		err error
	}{
		{
			schema: TFSchema{
				Name: "foo"	,
				SwaggerSpec: "spec1",
				PropertyLinks: map[string][]SwaggerLink{
					"p1": {
						{
							Spec: strPtr("spec2"),
							SchemaProp: *propertyaddr.NewPropertyAddrFromString("schema1:p1"),
						},
					},
				},
			},
			err: nil,
		},
		{
			schema: TFSchema{
				Name: "foo"	,
				SwaggerSpec: "spec1",
				PropertyLinks: map[string][]SwaggerLink{
					"foo:p1": {
						{
							Spec: strPtr("spec2"),
							SchemaProp: *propertyaddr.NewPropertyAddrFromString("schema1:p1"),
						},
					},
				},
			},
			err: errors.New("terraform property addr foo:p1 should not specify owner"),
		},
		{
			schema: TFSchema{
				Name: "foo"	,
				SwaggerSpec: "spec1",
				PropertyLinks: map[string][]SwaggerLink{
					"p1": {
						{
							Spec: strPtr("spec2"),
							SchemaProp: *propertyaddr.NewPropertyAddrFromString("p1.p2"),
						},
					},
				},
			},
			err: errors.New("swagger property addr p1.p2 should specify owner"),
		},
	}

	for idx, c :=  range cases {
		err :=c.schema.Validate()
		if c.err != nil {
			assert.EqualError(t, err, c.err.Error(), idx)
		} else {
			assert.NoError(t, err, idx)
		}
	}
}

func strPtr(s string) *string {
	return &s
}

