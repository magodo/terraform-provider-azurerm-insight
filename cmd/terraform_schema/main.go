package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/magodo/terraform-provider-azurerm-insight/pkg"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	schemaPath := flag.String("schema", "", "The relative path to the pkg provider schema file")
	outputPath := flag.String("output", pwd, "The output directory")
	resource := flag.String("resource", "", "The pkg resource to generate flattened schema. If not specified, will apply to all resources available.")
	isDataSource := flag.Bool("data-source", false, "Whether applies to data source")
	showHelp := flag.Bool("help", false, "Display this message")

	flag.Parse()

	if *showHelp {
		flag.Usage()
		return
	}

	// Read the provider schema
	b, err := ioutil.ReadFile(*schemaPath)
	if err != nil {
		log.Fatal(err)
	}

	var provider pkg.TerraformProvider
	if err := json.Unmarshal(b, &provider); err != nil {
		log.Fatal(err)
	}

	// Prepare output directory
	if err := os.MkdirAll(*outputPath, 0755); err != nil {
		log.Fatal(err)
	}

	if *resource != "" {
		var (
			ok     bool
			prefix string
			schema *pkg.TerraformSchema
		)
		if *isDataSource {
			schema, ok = provider.DataSourceSchemas[*resource]
			if !ok {
				log.Fatalf("No such data source: %s", *resource)
			}
			prefix = "data_"
		} else {
			schema, ok = provider.ResourceSchemas[*resource]
			if !ok {
				log.Fatalf("No such resource: %s", *resource)
			}
		}
		ofile := filepath.Join(*outputPath, prefix+*resource)
		if err := genFile(prefix+*resource, schema.Block, ofile); err != nil {
			log.Fatal(err)
		}
		return
	}

	var (
		schemas map[string]*pkg.TerraformSchema
		oprefix string
	)
	if *isDataSource {
		schemas = provider.DataSourceSchemas
		oprefix = "data_"
	} else {
		schemas = provider.ResourceSchemas
	}

	for res, schema := range schemas {
		ofile := filepath.Join(*outputPath, oprefix+res)
		if err := genFile(oprefix+res, schema.Block, ofile); err != nil {
			log.Fatal(err)
		}
	}
	return
}

func genFile(schemaName string, blk *pkg.TerraformBlock, ofileBase string) error {
	schema := pkg.NewSchemaScaffoldFromTerraformBlock(schemaName, blk)
	b, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return err
	}

	ofile := ofileBase + ".json"
	ofileBkp := ofileBase + ".json.bkp"

	// backup file if exists
	if stat, err := os.Stat(ofile); err == nil && stat.Mode().IsRegular() {
		src, err := os.Open(ofile)
		if err != nil {
			return err
		}

		dst, err := os.Create(ofileBkp)
		if err != nil {
			return err
		}

		if _, err := io.Copy(dst, src); err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(ofile, b, 0644); err != nil {
		return err
	}
	return nil

}
