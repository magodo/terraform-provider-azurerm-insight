package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/magodo/terraform-provider-azurerm-insight/pkg/core"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
	swaggerBaseDir := flag.String("swagger-base-dir", "", "The path to the swagger base directory (e.g. https://raw.githubusercontent.com/Azure/azure-rest-api-specs/master/specification)")
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

	swgschemas := core.NewSGWSchemas()
	err = filepath.Walk(*tfSchemaDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		var tfschema core.TFSchema
		if err := json.Unmarshal(b, &tfschema); err != nil {
			return err
		}

		if err := tfschema.LinkSwagger(swgschemas, *swaggerBaseDir); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatalf("error walking the terraform schema directory %q: %v\n", *tfSchemaDir, err)
	}

	b, err := json.MarshalIndent(swgschemas.GetAll(), "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(*outputPath, b, 0644); err != nil {
		log.Fatal(err)
	}
}
