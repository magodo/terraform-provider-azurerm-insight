package main

import (
	"fmt"

	"github.com/rivo/tview"
)

func PageSwaggerSchema(nextPage func()) (title string, content tview.Primitive) {
	infoView := tview.NewTextView()
	infoView.SetTitle("Coverage Information").SetBorder(true)

	var schemas []string
	for i := 0; i < 100; i++ {
		schemas = append(schemas, fmt.Sprintf("schema-%d", i))
	}
	schemaList := tview.NewList().ShowSecondaryText(false)
	schemaList.SetTitle("schema")
	schemaList.SetBorder(true)
	for _, schema := range schemas {
		schema := schema
		schemaList.AddItem(schema, "", 0, func() { fmt.Fprintf(infoView, schema) })
	}
	schemaList.SetBorderPadding(1, 1, 2, 2)

	var properties []string
	for i := 0; i < 100; i++ {
		properties = append(properties, fmt.Sprintf("propertie-%d", i))
	}
	propertieList := tview.NewList().ShowSecondaryText(false)
	propertieList.SetTitle("propertie")
	propertieList.SetBorder(true)
	for _, property := range properties {
		property := property
		propertieList.AddItem(property, "", 0, func() { fmt.Fprintf(infoView, property) })
	}
	propertieList.SetBorderPadding(1, 1, 2, 2)

	return "Table", tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(schemaList, 0, 1, true).
		AddItem(propertieList, 0, 1, true).
		AddItem(infoView, 0, 3, true)
}
