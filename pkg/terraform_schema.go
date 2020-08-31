package pkg

type TerraformProvider struct {
	ResourceSchemas   map[string]*TerraformSchema `json:"resource_schemas,omitempty"`
	DataSourceSchemas map[string]*TerraformSchema `json:"data_source_schemas,omitempty"`
}

type TerraformSchema struct {
	Block   *TerraformBlock `json:"block,omitempty"`
}

type TerraformBlock struct {
	Attributes map[string]*TerraformAttribute   `json:"attributes"`
	BlockTypes map[string]*TerraformNestedBlock `json:"block_types"`
}

type TerraformAttribute struct{}

type TerraformNestedBlock struct {
	TerraformBlock `json:"block"`
}
