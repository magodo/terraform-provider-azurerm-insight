package pkg

import (
	"encoding/json"
	"reflect"
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

	if !reflect.DeepEqual(schema, expect) {
		t.Fatalf(`
expect:

%+v

got:

%+v`, *expect, *schema)
	}
}

func TestMarshalTFSchema(t *testing.T) {
	tfschema := TFSchema{
		Name:        "res1",
		SwaggerSpec: "spec1",
		PropertyLinks: map[string][]SwaggerLink{
			"bar": {
				{
					Spec:       strPtr("xxx"),
					SchemaProp: propertyAddr{owner: "schema1", addrs: []string{"p1", "p2"}},
				},
				{
					Spec:       strPtr("yyy"),
					SchemaProp: propertyAddr{owner: "schema2", addrs: []string{"p3", "p4"}},
				},
			},
			"block_a::block_a_a::bar": {
				{
					Spec:       strPtr("xxx"),
					SchemaProp: propertyAddr{owner: "schema1", addrs: []string{"p1", "p2"}},
				},
				{
					Spec:       strPtr("yyy"),
					SchemaProp: propertyAddr{owner: "schema2", addrs: []string{"p3", "p4"}},
				},
			},
		},
	}

	expect := []byte(`{
    "Name": "res1",
    "PropertyLinks": {
        "bar": [
            {
                "prop": {
                    "addr": "p1.p2",
                    "owner": "schema1"
                },
                "spec": "xxx"
            },
            {
                "prop": {
                    "addr": "p3.p4",
                    "owner": "schema2"
                },
                "spec": "yyy"
            }
        ],
        "block_a::block_a_a::bar": [
            {
                "prop": {
                    "addr": "p1.p2",
                    "owner": "schema1"
                },
                "spec": "xxx"
            },
            {
                "prop": {
                    "addr": "p3.p4",
                    "owner": "schema2"
                },
                "spec": "yyy"
            }
        ]
    },
    "spec": "spec1"
}`)

	actual, err := json.Marshal(tfschema)
	if err != nil {
		t.Fatal(err)
	}

	if !jsonDeepEqual(t, actual, expect) {
		t.Fatalf(`
expect:

%+v

got:

%+v`, string(expect), string(actual))
	}
}

func TestUnmarshalTFSchema(t *testing.T) {
	jsonInput := []byte(`{
    "Name": "res1",
    "PropertyLinks": {
        "bar": [
            {
                "prop": {
                    "addr": "p1.p2",
                    "owner": "schema1"
                },
                "spec": "xxx"
            },
            {
                "prop": {
                    "addr": "p3.p4",
                    "owner": "schema2"
                },
                "spec": "yyy"
            }
        ],
        "block_a::block_a_a::bar": [
            {
                "prop": {
                    "addr": "p1.p2",
                    "owner": "schema1"
                },
                "spec": "xxx"
            },
            {
                "prop": {
                    "addr": "p3.p4",
                    "owner": "schema2"
                },
                "spec": "yyy"
            }
        ]
    },
    "spec": "spec1"
}`)

	expect := TFSchema{
		Name:        "res1",
		SwaggerSpec: "spec1",
		PropertyLinks: map[string][]SwaggerLink{
			"bar": {
				{
					Spec:       strPtr("xxx"),
					SchemaProp: propertyAddr{owner: "schema1", addrs: []string{"p1", "p2"}},
				},
				{
					Spec:       strPtr("yyy"),
					SchemaProp: propertyAddr{owner: "schema2", addrs: []string{"p3", "p4"}},
				},
			},
			"block_a::block_a_a::bar": {
				{
					Spec:       strPtr("xxx"),
					SchemaProp: propertyAddr{owner: "schema1", addrs: []string{"p1", "p2"}},
				},
				{
					Spec:       strPtr("yyy"),
					SchemaProp: propertyAddr{owner: "schema2", addrs: []string{"p3", "p4"}},
				},
			},
		},
	}

	var schema TFSchema
	if err := json.Unmarshal(jsonInput, &schema); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(schema, expect) {
		t.Fatalf(`
expect:

%+v

got:

%+v`, expect, schema)
	}
}

func strPtr(s string) *string {
	return &s
}

func jsonDeepEqual(t *testing.T, x, y []byte) bool {
	return reflect.DeepEqual(jsonNormalize(t, x), jsonNormalize(t, y))
}

func jsonNormalize(t *testing.T, in []byte) []byte {
	var tmp interface{}
	if err := json.Unmarshal(in, &tmp); err != nil {
		t.Fatal(err)
	}
	out, _ := json.Marshal(tmp)
	return out
}
