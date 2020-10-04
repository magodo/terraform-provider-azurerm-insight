package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
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
	swaggerSpecPath := flag.String("swagger-spec-path", "", "The path to the swagger spec directory, either a HTTP URI or local path (e.g. https://raw.githubusercontent.com/Azure/azure-rest-api-specs/master/specification)")
	showHelp := flag.Bool("help", false, "Display this message")
	githubToken := flag.String("github-token", "", "Github access token used to interact with github repos")

	flag.Parse()

	if *showHelp {
		flag.Usage()
		return
	}

	swgschemas, err := core.NewSWGSchemasFromTerraformSchema(*swaggerSpecPath, *tfSchemaDir, *swaggerGrantBaseDir)
	if err != nil {
		log.Fatal(err)
	}

	azureswgschemas := NewSWGResourceProviders(*swgschemas)

	swaggerURL, err := url.Parse(*swaggerSpecPath)
	if err != nil {
		log.Println(err)
	}
	if swaggerURL.Scheme == "http" || swaggerURL.Scheme == "https" {
		if err := azureswgschemas.CompleteSWGResourceProvidersViaGithubAPI(context.TODO(), &ghwalk.WalkOptions{Token: *githubToken, Reverse: true}); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := azureswgschemas.CompleteSWGResourceProvidersViaLocalFS(*swaggerSpecPath); err != nil {
			log.Fatal(err)
		}
	}

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
