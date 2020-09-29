package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/magodo/terraform-provider-azurerm-insight/pkg/core"

	"github.com/magodo/terraform-provider-azurerm-insight/pkg/core/propertyaddr"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type pageSwaggerItems struct {
	rpList         *tview.List
	apiList        *tview.List
	schemaList     *tview.List
	propertyTree   *tview.TreeView
	propertyDetail *tview.TextView
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
				items.propertyDetail.Clear()
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
		items.apiList.AddItem(k, "", 0,
			func() {
				refreshSchemaList(items, v.Schemas)
				app.SetFocus(items.schemaList)
				items.propertyDetail.Clear()
			})
	}
}

func refreshSchemaList(items pageSwaggerItems, swgschemas SWGSchemas) {
	items.schemaList.Clear()

	schemas := make([]string, 0, len(swgschemas))
	for k := range swgschemas {
		schemas = append(schemas, k)
	}
	sort.Strings(schemas)

	for _, k := range schemas {
		v := swgschemas[k]
		propCovered, propTotal := v.SchemaCoverage()
		items.schemaList.AddItem(k, fmt.Sprintf("[cov. %.2f%%]", 100*float64(propCovered)/float64(propTotal)), 0,
			func() {
				refreshPropertyTree(items, v)
				app.SetFocus(items.propertyTree)
				items.propertyDetail.Clear()
			})
	}
}

func refreshPropertyTree(items pageSwaggerItems, swgschema SWGSchema) {
	root := tview.NewTreeNode(".")
	items.propertyTree.SetRoot(root).SetCurrentNode(root)

	var (
		colorCoveredLeaf    = tcell.ColorWhite
		colorNotCoveredLeaf = tcell.NewHexColor(0xff6666)
		colorGrantedLeaf    = tcell.ColorDimGrey
		colorNonLeaf        = tcell.NewHexColor(0x3399ff)
	)

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
					cnode.SetColor(colorGrantedLeaf)
				} else if len(prop.TFLinks) == 0 {
					cnode.SetColor(colorNotCoveredLeaf)
				} else {
					cnode.SetColor(colorCoveredLeaf)
				}
			} else {
				cnode.SetColor(colorNonLeaf)
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

	items.propertyTree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) != 0 {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
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

	items.propertyTree.SetRoot(root)
}

func PageSwagger(swgrps SWGResourceProviders) tview.Primitive {
	rpList := tview.NewList().ShowSecondaryText(true)
	rpList.SetTitle("Resource Provider").SetBorder(true)

	apiList := tview.NewList().ShowSecondaryText(true)
	apiList.SetTitle("API Version").SetBorder(true)
	apiList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'h' {
			app.SetFocus(rpList)
			return nil
		}
		return event
	})

	schemaList := tview.NewList().ShowSecondaryText(true)
	schemaList.SetTitle("Schema").SetBorder(true)
	schemaList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'h' {
			app.SetFocus(apiList)
			return nil
		}
		return event
	})

	propertyTree := tview.NewTreeView()
	propertyTree.SetTitle("Property").SetBorder(true)
	propertyTree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'h' {
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
		AddItem(items.schemaList, 0, 1, true).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(items.propertyTree, 0, 3, true).
				AddItem(propertyDetail, 0, 1, true),
			0, 5, true,
		)
}
