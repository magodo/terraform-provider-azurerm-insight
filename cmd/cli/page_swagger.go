package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/magodo/terraform-provider-azurerm-insight/pkg/core"

	"github.com/magodo/terraform-provider-azurerm-insight/pkg/core/propertyaddr"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var (
	colorCoveredProperty    = tcell.ColorWhite
	colorNotCoveredProperty = tcell.NewHexColor(0xff6666)
	colorGrantedProperty    = tcell.ColorDimGrey
	colorObjectProperty     = tcell.NewHexColor(0x3399ff)
)

var (
	colorTextCoveredSchema    = "white"
	colorTextNotCoveredSchema = "red"
	colorTextGrantedSchema    = "grey"
)

type pageSwaggerItems struct {
	rpList         *tview.List
	apiList        *tview.List
	schemaList     *tview.List
	propertyTree   *tview.TreeView
	propertyDetail *tview.TextView
}

func drawProgressBar(percentage float64) string {
	return fmt.Sprintf("[%s%s] - %.2f%%", strings.Repeat("#", int(percentage*10)), strings.Repeat(" ", 10-int(percentage*10)), 100*percentage)
}

func refreshResourceProviderList(items pageSwaggerItems, swgrps SWGResourceProviders) {
	items.rpList.Clear()

	rps := make([]string, 0, len(swgrps))
	for k := range swgrps {
		rps = append(rps, k)
	}
	sort.Strings(rps)

	for _, k := range rps {
		v := swgrps[k]
		items.rpList.AddItem(k, "", 0,
			func() {
				refreshApiVersionList(items, v.Apis)
				app.SetFocus(items.apiList)
			})
	}
}

func refreshApiVersionList(items pageSwaggerItems, swgapis SWGResourceProviderAPIs) {
	items.apiList.Clear()

	apis := make([]string, 0, len(swgapis))
	for k := range swgapis {
		apis = append(apis, k)
	}
	sort.Strings(apis)

	for _, k := range apis {
		v := swgapis[k]

		var covered, total int

		swgschemas := v.Schemas
		for _, v := range swgschemas {
			if v.IsGranted {
				continue
			}
			total++
			propCovered, _ := v.SchemaCoverage()
			if propCovered != 0 {
				covered++
			}
		}

		var cov float64
		if total == 0 {
			cov = 0
		} else {
			cov = float64(covered) / float64(total)
		}

		items.apiList.AddItem(k, drawProgressBar(cov), 0,
			func() {
				refreshSchemaList(items, v.Schemas)
				app.SetFocus(items.schemaList)
			})
	}
}

func refreshSchemaList(items pageSwaggerItems, swgschemas SWGSchemas) {
	items.schemaList.Clear()

	// There is a bug (potentially in tview) which will cause the "swgschemas" in SetChangedFunc closure pointing to the "wrong" one
	// when switching RP. This causes panic since the schema specified in the mainText doesn't exist in the swgschemas.
	// Hence, we explicitly set the changed func to be nil so that the tview framework will not incorrectly call the "old" handler with "new" item.
	items.schemaList.SetChangedFunc(nil)

	schemas := make([]string, 0, len(swgschemas))
	for k := range swgschemas {
		schemas = append(schemas, k)
	}
	sort.Strings(schemas)

	var formatText = func(colorText, rawText string) string {
		return fmt.Sprintf("[%s]%s", colorText, rawText)
	}
	var parseText = func(mainText string) (colorText, rawText string) {
		matches := regexp.MustCompile(`\[(.+)](.+)`).FindStringSubmatch(mainText)
		return matches[1], matches[2]
	}
	for _, k := range schemas {
		v := swgschemas[k]
		propCovered, propTotal := v.SchemaCoverage()

		var (
			mainText      string
			secondaryText string
		)
		if v.IsGranted {
			mainText = formatText(colorTextGrantedSchema, k)
		} else if propCovered == 0 {
			mainText = formatText(colorTextNotCoveredSchema, k)
		} else {
			var cov float64
			if propTotal == 0 {
				cov = 0
			} else {
				cov = float64(propCovered) / float64(propTotal)
			}
			mainText = formatText(colorTextCoveredSchema, k)
			secondaryText = fmt.Sprintf("%s\n", drawProgressBar(cov))
		}

		items.schemaList.AddItem(mainText, secondaryText, 0, nil)

		items.schemaList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			_, rawMainText := parseText(mainText)
			v := swgschemas[rawMainText]
			refreshPropertyTree(items, *v)
		})

		items.schemaList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			_, rawMainText := parseText(mainText)
			v := swgschemas[rawMainText]
			refreshPropertyTree(items, *v)
			app.SetFocus(items.propertyTree)
		})
	}
}

func refreshPropertyTree(items pageSwaggerItems, swgschema SWGSchema) {
	root := tview.NewTreeNode(".")
	items.propertyTree.SetRoot(root).SetCurrentNode(root)

	propaddrs := make([]string, 0, len(swgschema.Properties))
	for k := range swgschema.Properties {
		propaddrs = append(propaddrs, k)
	}
	sort.Strings(propaddrs)

	for _, addr := range propaddrs {
		prop := swgschema.Properties[addr]
		node := root

		curaddr := *propertyaddr.NewPropertyAddrFromString("")
		addrs := propertyaddr.NewPropertyAddrFromString(addr).RelativeAddrs()
		for idx, segment := range addrs {
			var cnode *tview.TreeNode
			for _, c := range node.GetChildren() {
				if c.GetText() == segment {
					cnode = c
					break
				}
			}
			isLeaf := idx == len(addrs)-1

			if cnode == nil {
				nodeText := segment
				cnode = tview.NewTreeNode(nodeText).SetExpanded(false)
				node.AddChild(cnode)
			}

			if isLeaf {
				cnode.SetReference(*prop)
				if prop.IsGranted {
					cnode.SetColor(colorGrantedProperty)
				} else if len(prop.TFLinks) == 0 {
					cnode.SetColor(colorNotCoveredProperty)
				} else {
					cnode.SetColor(colorCoveredProperty)
				}
			} else {
				cnode.SetColor(colorObjectProperty)
				curaddr = curaddr.Append(segment)
				cov, total, ok := swgschema.FindCoverage(curaddr)
				var coverage float64
				if ok {
					coverage = float64(cov) / float64(total)
				}
				cnode.SetReference(coverage)
			}
			node = cnode
		}
	}

	items.propertyTree.SetChangedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}

		showNodeInfo := func(cov float64) {
			items.propertyDetail.Clear()

			fmt.Fprintf(items.propertyDetail, "Coverage: %.2f%%", 100*cov)
		}
		showLeafNodeInfo := func(prop core.SWGSchemaProperty) {
			items.propertyDetail.Clear()

			if prop.IsGranted {
				fmt.Fprintf(items.propertyDetail, "Deliberately not supported in Terraform: %s", prop.GrantComment)
				return
			}
			tfproperties := make([]string, 0, len(prop.TFLinks))
			for _, tflink := range prop.TFLinks {
				prop := tflink.Prop
				tfproperties = append(tfproperties, fmt.Sprintf("- %s: %s", prop.Owner(), prop.RelativeAddrs().String()))
			}

			if len(tfproperties) == 0 {
				fmt.Fprintf(items.propertyDetail, "To be supported in Terraform in the future.")
			} else {
				fmt.Fprintf(items.propertyDetail, `Related Terraform Properties:

%s
`, strings.Join(tfproperties, "\n"))
			}
		}

		switch ref := reference.(type) {
		case float64:
			showNodeInfo(ref)
		case core.SWGSchemaProperty:
			showLeafNodeInfo(ref)
		}
	})

	items.propertyTree.SetSelectedFunc(func(node *tview.TreeNode) {
		children := node.GetChildren()
		if len(children) != 0 {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

	items.propertyTree.SetRoot(root)
}

func PageSwagger(swgrps SWGResourceProviders) tview.Primitive {
	rpList := tview.NewList().ShowSecondaryText(true)
	rpList.SetTitle("Resource Provider").SetBorder(true)

	apiList := tview.NewList().ShowSecondaryText(true)
	apiList.SetTitle("API Version").SetBorder(true)
	apiList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.SetFocus(rpList)
			return nil
		}
		return event
	})

	schemaList := tview.NewList().ShowSecondaryText(true)
	schemaList.SetTitle("Schema").SetBorder(true)
	schemaList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.SetFocus(apiList)
			return nil
		}
		return event
	})

	propertyTree := tview.NewTreeView()
	propertyTree.SetTitle("Property").SetBorder(true)
	propertyTree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.SetFocus(schemaList)
			return nil
		}
		return event
	})
	propertyDetail := tview.NewTextView()
	propertyDetail.SetTitle("PropertyDetail").SetBorder(true)

	items := pageSwaggerItems{
		rpList:         rpList,
		apiList:        apiList,
		schemaList:     schemaList,
		propertyTree:   propertyTree,
		propertyDetail: propertyDetail,
	}

	refreshResourceProviderList(items, swgrps)

	return tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(items.rpList, 0, 1, true).
		AddItem(items.apiList, 0, 1, true).
		AddItem(items.schemaList, 0, 2, true).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(items.propertyTree, 0, 3, true).
				AddItem(propertyDetail, 0, 1, true),
			0, 5, true,
		)
}
