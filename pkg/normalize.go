package pkg

import (
	"fmt"
	openapispec "github.com/go-openapi/spec"
	"net/url"
	"path/filepath"
	"strings"
)

// normalize URI reference against the base path, for local files only.
func normalizeFileRef(ref *openapispec.Ref, relativeBase string) *openapispec.Ref {
	if ref.String() == "" {
		r, _ := openapispec.NewRef(relativeBase)
		return &r
	}
	s := normalizePaths(ref.String(), relativeBase)
	r, _ := openapispec.NewRef(s)
	return &r
}

// normalize URI reference string against the base path, for local files only.
func normalizePaths(refPath, base string) string {
	if filepath.IsAbs(refPath) {
		parts := strings.Split(refPath, "#")
		result := filepath.FromSlash(parts[0])
		if len(parts) == 2 {
			result += "#" + parts[1]
		}
		return result
	}

	refURL, _ := url.Parse(refPath)
	baseURL, _ := url.Parse(base)
	if !strings.HasPrefix(refPath, "#") {
		// combining paths
		newBase := fmt.Sprintf("%s#%s", filepath.Join(filepath.Dir(base), filepath.FromSlash(refURL.Path)), refURL.Fragment)
		return newBase

	}
	// copying fragment from ref to base
	baseURL.Fragment = refURL.Fragment
	return baseURL.String()
}
