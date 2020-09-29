package core

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type SWGGrant map[SWGSchemaAddr]SWGSchemaGrant

type SWGSchemaGrant struct {
	// The grant comment why the schema is granted.
	Comment string `json:",omitempty"`

	// Property grant map, whose key is the propertyaddr, whose value is the grant comment.
	Properties map[string]string `json:",omitempty"`
}

func (g SWGSchemaGrant) IsSchemaGranted() bool {
	return len(g.Properties) == 0
}

// NewSWGGrantFromFiles construct a SWGGrant from a grantBaseDir which contains the
// folder layout as defined by the SWGSchemaAddr.
func NewSWGGrantFromFiles(grantBaseDir string) (SWGGrant, error) {
	var swgGrant SWGGrant = map[SWGSchemaAddr]SWGSchemaGrant{}
	return swgGrant, filepath.Walk(grantBaseDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() ||
			(!strings.HasSuffix(info.Name(), ".json") &&
				!strings.HasSuffix(info.Name(), ".yaml") &&
				!strings.HasSuffix(info.Name(), ".yml")) {
			return err
		}

		infileSwgGrant := map[string]SWGSchemaGrant{}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(b, &infileSwgGrant); err != nil {
			return err
		}

		relPath, err := filepath.Rel( grantBaseDir, path)
		if err != nil {
			return err
		}

		for schemaName, schemaGrant := range infileSwgGrant {
			swgGrant[NewSWGSchemaAddr(relPath, schemaName)] = schemaGrant
		}
		return nil
	})
}
