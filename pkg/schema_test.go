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

	var expect Schema = map[string][]Link{
		"bar":                     nil,
		"block_a::block_a_a::bar": nil,
		"block_a::foo":            nil,
		"foo":                     nil,
	}

	var block TerraformBlock
	if err := json.Unmarshal(jsonInput, &block); err != nil {
		t.Fatal(err)
	}
	schema, err := NewSchemaFromTerraformBlock(&block)
	if err != nil {
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

func TestUnmarshalSchema(t *testing.T) {
	jsonInput := []byte(`
{
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
`)

	var expect Schema = map[string][]Link{
		"bar": {
			{
				SwaggerSpecPath:   "xxx",
				SwaggerSchemaProp: SwaggerSchemaProp{segments: []string{"a", "b"}},
			},
			{
				SwaggerSpecPath:   "yyy",
				SwaggerSchemaProp: SwaggerSchemaProp{segments: []string{"b", "c"}},
			},
		},
		"block_a::block_a_a::bar": {
			{
				SwaggerSpecPath:   "xxx",
				SwaggerSchemaProp: SwaggerSchemaProp{segments: []string{"a", "b"}},
			},
			{
				SwaggerSpecPath:   "yyy",
				SwaggerSchemaProp: SwaggerSchemaProp{segments: []string{"b", "c"}},
			},
		},
	}

	var schema Schema
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
