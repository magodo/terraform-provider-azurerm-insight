package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"encoding/json"

	"github.com/magodo/terraform-provider-azurerm-insight/pkg"
)

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Shed insight on the Terraform AzureRM Provider.\n\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	queryCmd := flag.NewFlagSet("query", flag.ExitOnError)
	swaggerRelPath := queryCmd.String("swagger-path", "", "Swagger file relative path")
	schemaName := queryCmd.String("schema-name", "", "Name of the Swagger schema")

	if len(os.Args) < 2 {
		flag.Usage()
	}

	switch os.Args[1] {
	case "help":
		flag.Usage()
	case "query":
		queryCmd.Parse(os.Args[2:])
		if *swaggerRelPath == "" || *schemaName == "" {
			queryCmd.Usage()
			os.Exit(2)
		}
		runQueryCmd(*swaggerRelPath, *schemaName)
	default:
		flag.Usage()
	}

}

func mustGetSwgSchemas() pkg.SWGSchemas {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	swgschemaPath := filepath.Join(pwd, "swagger_schema.json")
	b, err := ioutil.ReadFile(swgschemaPath)
	if err != nil {
		log.Fatal(err)
	}

	var swgschemas pkg.SWGSchemas
	if err := json.Unmarshal(b, &swgschemas); err != nil {
		log.Fatal(err)
	}
	return swgschemas
}

func runQueryCmd(swaggerRelPath, schemaName string) {
}
