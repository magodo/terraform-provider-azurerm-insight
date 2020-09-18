package pkg

import (
	"fmt"
	openapispec "github.com/go-openapi/spec"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

// relativeBase could be an ABSOLUTE file path or an ABSOLUTE URL
func normalizeFileRef(ref *openapispec.Ref, relativeBase string) *openapispec.Ref {
	if ref.String() == "" {
		r, _ := openapispec.NewRef(relativeBase)
		return &r
	}

	s := normalizePaths(ref.String(), relativeBase)
	r, _ := openapispec.NewRef(s)
	return &r
}

// base or refPath could be a file path or a URL
// given a base absolute path and a ref path, return the absolute path of refPath
// 1) if refPath is absolute, return it
// 2) if refPath is relative, join it with basePath keeping the scheme, hosts, and ports if exists
// base could be a directory or a full file path
func normalizePaths(refPath, base string) string {
	refURL, _ := url.Parse(refPath)
	if path.IsAbs(refURL.Path) || filepath.IsAbs(refPath) {
		// refPath is actually absolute
		if refURL.Host != "" {
			return refPath
		}
		parts := strings.Split(refPath, "#")
		result := filepath.FromSlash(parts[0])
		if len(parts) == 2 {
			result += "#" + parts[1]
		}
		return result
	}

	// relative refPath
	baseURL, _ := url.Parse(base)
	if !strings.HasPrefix(refPath, "#") {
		// combining paths
		if baseURL.Host != "" {
			baseURL.Path = path.Join(path.Dir(baseURL.Path), refURL.Path)
		} else { // base is a file
			newBase := fmt.Sprintf("%s#%s", filepath.Join(filepath.Dir(base), filepath.FromSlash(refURL.Path)), refURL.Fragment)
			return newBase
		}

	}
	// copying fragment from ref to base
	baseURL.Fragment = refURL.Fragment
	return baseURL.String()
}
