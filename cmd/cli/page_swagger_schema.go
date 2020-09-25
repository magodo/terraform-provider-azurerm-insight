package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func PageSwaggerSchema(nextPage func()) (title string, content tview.Primitive) {
	infoView := tview.NewTextView()
	infoView.SetTitle("Coverage Information").SetBorder(true)

	propertieList := tview.NewList().ShowSecondaryText(true)
	propertieList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'j' {
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		}
		return event
	})
	schemaList := tview.NewList().ShowSecondaryText(true)

	var schemas []string
	for i := 0; i < 100; i++ {
		schemas = append(schemas, fmt.Sprintf("schema-%d", i))
	}
	schemaList.SetTitle("schema").SetBorder(true)
	for _, schema := range schemas {
		schema := schema
		schemaList.AddItem(schema, "30%", 0,
			func() {
				app.SetFocus(propertieList)
				infoView.Clear()
				fmt.Fprintf(infoView, "schema info")
			})
	}
	schemaList.SetBorderPadding(1, 1, 1, 1)

	var properties []string
	for i := 0; i < 100; i++ {
		properties = append(properties, fmt.Sprintf("propertie-%d", i))
	}
	propertieList.SetTitle("propertie").SetBorder(true)
	for _, property := range properties {
		property := property
		propertieList.AddItem(property, "20%", 0, func() { fmt.Fprintf(infoView, property) })
	}
	propertieList.SetBorderPadding(1, 1, 1, 1)
	propertieList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.SetFocus(schemaList)
			return nil
		}
		return event
	})

	return "Table", tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(schemaList, 0, 1, true).
		AddItem(propertieList, 0, 1, true).
		AddItem(infoView, 0, 3, true)
}
