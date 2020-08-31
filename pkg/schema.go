package pkg

import (
	"strings"
)

const addrSep = "::"

type Addr struct {
	segments []string
}

type Schema map[string][]Link

func (addr Addr) String() string {
	return strings.Join(addr.segments, addrSep)
}

func (addr Addr) Append(oaddr string) Addr {
	segments := make([]string, len(addr.segments))
	copy(segments, addr.segments)
	segments = append(segments, oaddr)
	return Addr{segments: segments}
}

func (addr *Addr) Pop() string {
	if len(addr.segments) > 1 {
		addr.segments = addr.segments[:len(addr.segments)-1]
	}
	return addr.segments[len(addr.segments)-1]
}

func NewSchemaFromTerraformBlock(block *TerraformBlock) (*Schema, error) {
	var schema Schema = map[string][]Link{}
	recordAttributeWithinBlock(Addr{}, schema, block)
	return &schema, nil
}

func recordAttributeWithinBlock(parentBlockAddr Addr, schema Schema, block *TerraformBlock) {
	for attrKey := range block.Attributes {
		addr := parentBlockAddr.Append(attrKey)
		schema[addr.String()] = []Link{}
	}
	for blockKey, blockVal := range block.BlockTypes {
		addr := parentBlockAddr.Append(blockKey)
		recordAttributeWithinBlock(addr, schema, &blockVal.TerraformBlock)
	}
}
