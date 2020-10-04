package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/magodo/terraform-provider-azurerm-insight/pkg/core"
)

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Generate swagger schema file, containing complete links, based on the terraform schema.\n\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	tfSchemaDir := flag.String("tf-schema-dir", "", "The path to the directory contains terraform schemas")
	swaggerGrantBaseDir := flag.String("swagger-grant-dir", "", "The path to the base directory contains swagger grant info (e.g. azure_knowledgebase/swagger_grants)")
	swaggerSpecPath := flag.String("swagger-spec-path", "", "The path to the swagger spec directory, either a HTTP URI or local path (e.g. https://raw.githubusercontent.com/Azure/azure-rest-api-specs/master/specification)")
	outputPath := flag.String("output", filepath.Join(pwd, "swagger_schema.json"), "The output file")
	showHelp := flag.Bool("help", false, "Display this message")

	flag.Parse()

	if *showHelp {
		flag.Usage()
		return
	}

	// Prepare output directory
	outputDir := filepath.Dir(*outputPath)
	stat, err := os.Stat(outputDir)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.IsDir() {
		log.Fatalf("%q exists but is not a directory", outputDir)
	}

	os.Remove(outputDir)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatal(err)
	}

	swgschemas, err := core.NewSWGSchemasFromTerraformSchema(*swaggerSpecPath, *tfSchemaDir, *swaggerGrantBaseDir)
	if err != nil {
		log.Fatal(err)
	}

	// Construct a temporary type to include the property coverage info in schema level.
	type swgSchemaWithCoverage struct {
		Coverage float64
		*core.SWGSchema
	}
	schemaMap := map[core.SWGSchemaAddr]swgSchemaWithCoverage{}
	for schemaAddr, schema := range swgschemas.GetAll() {
		covered, total := schema.SchemaCoverage()
		schemaMap[schemaAddr] = swgSchemaWithCoverage{
			Coverage:  float64(covered) / float64(total),
			SWGSchema: schema,
		}
	}

	b, err := json.MarshalIndent(schemaMap, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(*outputPath, b, 0644); err != nil {
		log.Fatal(err)
	}
}
