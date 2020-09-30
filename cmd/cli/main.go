package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/magodo/ghwalk"

	"github.com/magodo/terraform-provider-azurerm-insight/pkg/core"
	"github.com/rivo/tview"
)

// The application.
var app = tview.NewApplication()

func main() {

	// Parse CLI flags
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "CLI Viewer of internal mapping between Terraform schema and Swagger schema for terraform-provider-azurerm.\n\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	tfSchemaDir := flag.String("tf-schema-dir", "", "The path to the directory contains terraform schemas")
	swaggerGrantBaseDir := flag.String("swagger-grant-dir", "", "The path to the base directory contains swagger grant info (e.g. azure_knowledgebase/swagger_grants)")
	swaggerBaseDir := flag.String("swagger-base-dir", "", "The path to the swagger base directory (e.g. https://raw.githubusercontent.com/Azure/azure-rest-api-specs/master/specification)")
	showHelp := flag.Bool("help", false, "Display this message")
	githubToken := flag.String("github-token", "", "Github access token used to interact with github repos (e.g. azure-rest-api-spec which holds the Azure Swagger Spec)")

	flag.Parse()

	if *showHelp {
		flag.Usage()
		return
	}

	swgschemas, err := core.NewSWGSchemasFromTerraformSchema(*swaggerBaseDir, *tfSchemaDir, *swaggerGrantBaseDir)
	if err != nil {
		log.Fatal(err)
	}

	azureswgschemas := NewSWGResourceProviders(*swgschemas)
	azureswgschemas.CompleteSWGResourceProviders(context.TODO(), &ghwalk.WalkOptions{Token: *githubToken, Reverse: true})

	page := PageSwagger(azureswgschemas)

	// Create the main layout.
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(page, 0, 1, true)

	// Start the application.
	if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
