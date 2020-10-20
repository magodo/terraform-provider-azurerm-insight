package core

import (
	"fmt"
	"testing"

	"github.com/magodo/terraform-provider-azurerm-insight/pkg/core/propertyaddr"
	"github.com/stretchr/testify/require"
)

func TestNewSWGPropertyCoverageStore_Add(t *testing.T) {
	type swgproperty struct {
		addr propertyaddr.SwaggerPropertyAddr
		prop SWGSchemaProperty
	}
	cases := []struct {
		propertiesToAdd []swgproperty
		expectStore     SWGPropertyCoverageStore
	}{
		{
			propertiesToAdd: []swgproperty{
				{
					addr: propertyaddr.MustParseSwaggerPropertyAddr("prop1.not_covered"),
					prop: SWGSchemaProperty{
						TFLinks: []TFLink{},
					},
				},
				{
					addr: propertyaddr.MustParseSwaggerPropertyAddr("prop1.covered"),
					prop: SWGSchemaProperty{
						TFLinks: []TFLink{{}},
					},
				},
				{
					addr: propertyaddr.MustParseSwaggerPropertyAddr("prop2.covered"),
					prop: SWGSchemaProperty{
						TFLinks: []TFLink{{}},
					},
				},
				{
					addr: propertyaddr.MustParseSwaggerPropertyAddr("prop_granted"),
					prop: SWGSchemaProperty{
						TFLinks:   []TFLink{},
						IsGranted: true,
					},
				},
			},
			expectStore: SWGPropertyCoverageStore{
				node: swgPropertyCoverageNode{
					TotalAmount:   3,
					CoveredAmount: 2,
					Children: map[string]*swgPropertyCoverageNode{
						"prop1": {
							TotalAmount:   2,
							CoveredAmount: 1,
							Children: map[string]*swgPropertyCoverageNode{
								"not_covered": {
									TotalAmount:   1,
									CoveredAmount: 0,
									Children:      map[string]*swgPropertyCoverageNode{},
								},
								"covered": {
									TotalAmount:   1,
									CoveredAmount: 1,
									Children:      map[string]*swgPropertyCoverageNode{},
								},
							},
						},
						"prop2": {
							TotalAmount:   1,
							CoveredAmount: 1,
							Children: map[string]*swgPropertyCoverageNode{
								"covered": {
									TotalAmount:   1,
									CoveredAmount: 1,
									Children:      map[string]*swgPropertyCoverageNode{},
								},
							},
						},
					},
				},
			},
		},
	}

	for idx, c := range cases {
		store := NewSWGPropertyCoverageStore()
		for idx2, prop := range c.propertiesToAdd {
			require.NoError(t, store.Add(prop.addr, prop.prop), fmt.Sprintf("%d.%d", idx, idx2))
		}
		require.Equal(t, c.expectStore, store, idx)
	}
}

func TestNewSWGPropertyCoverageStore_FindCoverage(t *testing.T) {
	type result struct {
		total   int
		covered int
		ok      bool
	}

	type subtest struct {
		propAddr string
		expect   result
	}

	cases := []struct {
		store               SWGPropertyCoverageStore
		subtest             []subtest
		expectTotalCoverage result
	}{
		{
			store: SWGPropertyCoverageStore{
				node: swgPropertyCoverageNode{
					TotalAmount:   3,
					CoveredAmount: 1,
					Children: map[string]*swgPropertyCoverageNode{
						"prop1": {
							TotalAmount:   2,
							CoveredAmount: 1,
							Children: map[string]*swgPropertyCoverageNode{
								"nest": {
									TotalAmount:   2,
									CoveredAmount: 1,
									Children: map[string]*swgPropertyCoverageNode{
										"covered": {
											TotalAmount:   1,
											CoveredAmount: 1,
											Children:      map[string]*swgPropertyCoverageNode{},
										},
										"uncovered": {
											TotalAmount:   1,
											CoveredAmount: 0,
											Children:      map[string]*swgPropertyCoverageNode{},
										},
									},
								},
							},
						},
						"uncovered": {
							TotalAmount:   1,
							CoveredAmount: 0,
							Children:      map[string]*swgPropertyCoverageNode{},
						},
					},
				},
			},
			subtest: []subtest{
				{
					propAddr: "non_existed",
					expect: result{
						ok: false,
					},
				},
				{
					propAddr: "uncovered",
					expect: result{
						total: 1,
						ok:    true,
					},
				},
				{
					propAddr: "prop1.nest.covered",
					expect: result{
						total:   1,
						covered: 1,
						ok:      true,
					},
				},
				{
					propAddr: "prop1.nest.uncovered",
					expect: result{
						total:   1,
						covered: 0,
						ok:      true,
					},
				},
				{
					propAddr: "prop1.nest",
					expect: result{
						total:   2,
						covered: 1,
						ok:      true,
					},
				},
				{
					propAddr: "prop1",
					expect: result{
						total:   2,
						covered: 1,
						ok:      true,
					},
				},
			},
			expectTotalCoverage: result{
				total:   3,
				covered: 1,
				ok:      true,
			},
		},
	}

	for idx, c := range cases {
		for idx2, subtest := range c.subtest {
			covered, total, ok := c.store.FindCoverage(propertyaddr.MustParseSwaggerPropertyAddr(subtest.propAddr))
			result := result{
				total:   total,
				covered: covered,
				ok:      ok,
			}
			require.Equal(t, subtest.expect, result, fmt.Sprintf("%d.%d", idx, idx2))
		}

		covered, total := c.store.SchemaCoverage()
		result := result{
			total:   total,
			covered: covered,
			ok:      true,
		}
		require.Equal(t, c.expectTotalCoverage, result, idx)
	}
}
