package core

import (
	"fmt"
	"github.com/magodo/terraform-provider-azurerm-insight/pkg/core/propertyaddr"
)

type SWGPropertyCoverageStore struct {
	node swgPropertyCoverageNode
}

func NewSWGPropertyCoverageStore() SWGPropertyCoverageStore {
	return SWGPropertyCoverageStore{
		node: swgPropertyCoverageNode{
			Children: map[string]*swgPropertyCoverageNode{},
		},
	}
}

// Add adds a SWGSchemaProperty and record the coverage state in each property level.
// If a property is added more than once, an error will be returned.
// NOTE: The granted property will be ignored
func (store *SWGPropertyCoverageStore) Add(propAddr propertyaddr.PropertyAddr, prop SWGSchemaProperty) error {
	if prop.IsGranted {
		return nil
	}

	addrs := propAddr.RelativeAddrs()
	isCovered := len(prop.TFLinks) != 0

	if _, _, ok := store.FindCoverage(propAddr); ok {
		return fmt.Errorf("%q has already been added to the coverage store", propAddr)
	}

	if len(addrs) == 0 {
		return nil
	}

	store.node.add(addrs, isCovered)
	return nil
}

func (store *SWGPropertyCoverageStore) SchemaCoverage() (covered, total int) {
	return store.node.CoveredAmount, store.node.TotalAmount
}

func (store *SWGPropertyCoverageStore) FindCoverage(propAddr propertyaddr.PropertyAddr) (covered, total int, ok bool) {
	addrs := propAddr.RelativeAddrs()
	node := store.node
	for _, addr := range addrs {
		tmpNode, ok := node.Children[addr]
		if !ok {
			return 0, 0, false
		}
		node = *tmpNode
		covered, total = node.CoveredAmount, node.TotalAmount
	}
	return covered, total, true
}

type swgPropertyCoverageNode struct {
	TotalAmount   int
	CoveredAmount int
	Children      map[string]*swgPropertyCoverageNode
}

func (node *swgPropertyCoverageNode) add(addrs propertyaddr.RelativeAddrs, isCovered bool) {
	node.TotalAmount++
	if isCovered {
		node.CoveredAmount++
	}

	if len(addrs) == 0 {
		return
	}

	child, ok := node.Children[addrs[0]]
	if !ok {
		node.Children[addrs[0]] = &swgPropertyCoverageNode{
			Children: map[string]*swgPropertyCoverageNode{},
		}
		child = node.Children[addrs[0]]
	}
	child.add(addrs[1:], isCovered)

	return
}
