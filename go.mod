module github.com/magodo/terraform-provider-azurerm-insight

go 1.15

require (
	github.com/gdamore/tcell v1.4.0
	github.com/go-openapi/loads v0.19.5
	github.com/go-openapi/spec v0.19.8
	github.com/magodo/ghwalk v0.0.0-20200930074045-b9a34d077a8b
	github.com/rivo/tview v0.0.0-20200915114512-42866ecf6ca6
	github.com/stretchr/testify v1.6.1
	github.com/zclconf/go-cty v1.6.1
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace github.com/go-openapi/spec => github.com/magodo/spec v0.19.10-0.20201124144715-3e5006560d1f
