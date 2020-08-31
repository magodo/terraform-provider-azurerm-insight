package pkg

import "strings"

type Link struct {
	SwaggerSpecPath   string            `json:"spec"` // Relative path to swagger spec file from the azure rest api spec repo root
	SwaggerSchemaProp SwaggerSchemaProp `json:"prop"` // dot-separated property, starting from the schema used as the PUT body parameter
}

type SwaggerSchemaProp struct {
	segments []string
}

func (p SwaggerSchemaProp) String() string {
	return strings.Join(p.segments, ".")
}

func (a *SwaggerSchemaProp) UnmarshalJSON(b []byte) error {
	a.segments = strings.Split(string(b), ".")
	return nil
}

func (a SwaggerSchemaProp) MarshalJSON() ([]byte, error) {
	return []byte(strings.Join(a.segments, ".")), nil
}
