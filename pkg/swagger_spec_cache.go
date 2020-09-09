package pkg

import (
	"fmt"
	"github.com/go-openapi/loads"
	"sync"
)

type SwaggerSpecCache struct {
	sync.Mutex
	m map[string]*loads.Document
}

// swaggerSpecCache caches the swagger document (spec) using swagger file absolute path as key.
var swaggerSpecCache SwaggerSpecCache

// LoadSwaggerSpec load a certain swagger spec (document)
func LoadSwaggerSpec(spec string) (*loads.Document, error) {
	swaggerSpecCache.Lock()
	defer swaggerSpecCache.Unlock()

	// construct key
	if schema, ok := swaggerSpecCache.m[spec]; ok {
		return schema, nil
	}

	doc, err := loads.Spec(spec)
	if err != nil {
		return nil, fmt.Errorf("loading swagger spec %s: %w", spec, err)
	}

	swaggerSpecCache.m[spec] = doc
	return doc, nil
}
