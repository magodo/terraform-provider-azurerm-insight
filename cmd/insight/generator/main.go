package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/magodo/terraform-provider-azurerm-insight/pkg"
)

// Generate the swgschema from tfschema
//go:generate go run ../../swagger_schema/main.go  -swagger-base-dir ../../../assets/static/azure-rest-api-specs/specification -tf-schema-dir ../../../assets/terraform_schema
func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Generate Go code from swg schema metadata and tf schema metadata.\n\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	tfSchemaDir := flag.String("tf-schema-dir", "", "The path to the directory contains terraform schema metadata")
	swgPath := flag.String("swg-schema", "", "The path to the swg schema metadata")
	outputDir := flag.String("output-dir", "", "The path to the directory where to generate the code")

	flag.Parse()

	if *tfSchemaDir == "" || *swgPath == "" || *outputDir == "" {
		flag.Usage()
	}

	var (
		swgschemas pkg.SWGSchemas
		tfschemas  map[string]*pkg.TFSchema
	)

	b, err := ioutil.ReadFile(*swgPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(b, &swgschemas); err != nil {
		log.Fatal(err)
	}

}
