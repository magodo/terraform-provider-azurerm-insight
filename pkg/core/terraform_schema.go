package core

import "github.com/zclconf/go-cty/cty"

type TerraformProviderSchemas struct {
	FormatVersion string                       `json:"format_version"`
	Schemas       map[string]TerraformProvider `json:"provider_schemas"`
}

type TerraformProvider struct {
	ResourceSchemas   map[string]*TerraformSchema `json:"resource_schemas,omitempty"`
	DataSourceSchemas map[string]*TerraformSchema `json:"data_source_schemas,omitempty"`
}

type TerraformSchema struct {
	Block *TerraformBlock `json:"block,omitempty"`
}

type TerraformBlock struct {
	Attributes map[string]*TerraformAttribute   `json:"attributes"`
	BlockTypes map[string]*TerraformNestedBlock `json:"block_types"`
}

type TerraformAttribute struct {
	Type *cty.Type `json:"type"`
}

type TerraformNestedBlock struct {
	TerraformBlock `json:"block"`
}
