package pkg

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewSchemaFromTerraformBlock(t *testing.T) {
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

	expect := TFSchema{
		Name: "res1",
		PropertyLinks: map[string][]SwaggerLink{
			"bar":                   nil,
			"block_a.block_a_a.bar": nil,
			"block_a.foo":           nil,
			"foo":                   nil,
		},
	}

	var block TerraformBlock
	if err := json.Unmarshal(jsonInput, &block); err != nil {
		t.Fatal(err)
	}
	schema := NewSchemaScaffoldFromTerraformBlock("res1", &block)

	if reflect.DeepEqual(schema, expect) {
		t.Fatalf(`
expect:

%+v

got:

%+v`, expect, schema)
	}
}

func TestUnmarshalSchema(t *testing.T) {
	jsonInput := []byte(`
{
    "name": "res1",
    "spec": "spec1",
    "attributes": {
		"bar": [
		  {
			"spec": "xxx",
			"prop": "a.b"
		  },
		  {
			"spec": "yyy",
			"prop": "b.c"
		  }
		],
		"block_a::block_a_a::bar": [
		  {
			"spec": "xxx",
			"prop": "a.b"
		  },
		  {
			"spec": "yyy",
			"prop": "b.c"
		  }
		]
    }
}
`)

	expect := TFSchema{
		Name:        "res1",
		SwaggerSpec: "spec1",
		PropertyLinks: map[string][]SwaggerLink{
			"bar": {
				{
					Spec:              "xxx",
					SwaggerSchemaProp: propertyAddr{segments: []string{"a", "b"}},
				},
				{
					Spec:              "yyy",
					SwaggerSchemaProp: propertyAddr{segments: []string{"b", "c"}},
				},
			},
			"block_a::block_a_a::bar": {
				{
					Spec:              "xxx",
					SwaggerSchemaProp: propertyAddr{segments: []string{"a", "b"}},
				},
				{
					Spec:              "yyy",
					SwaggerSchemaProp: propertyAddr{segments: []string{"b", "c"}},
				},
			},
		},
	}

	var schema TFSchema
	if err := json.Unmarshal(jsonInput, &schema); err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(schema, expect) {
		t.Fatalf(`
expect:

%+v

got:

%+v`, expect, schema)
	}
}
