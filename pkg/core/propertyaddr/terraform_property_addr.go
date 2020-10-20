package propertyaddr

import (
	"encoding/json"
	"strings"
)

const terraformPropertyResourceSep = ":"
const terraformPropertyAddrSep = "."

type TerraformPropertyAddr struct {
	ResourceName string
	PropertyAddr TerraformRelativeAddrs
}

type TerraformRelativeAddrs []string

func (a TerraformRelativeAddrs) String() string {
	return strings.Join(a, terraformPropertyAddrSep)
}

func ParseTerraformPropertyAddr(addr string) *TerraformPropertyAddr {
	p := strings.Split(addr, terraformPropertyResourceSep)
	switch len(p) {
	case 1:
		addrs := []string{}
		if p[0] != "" {
			addrs = strings.Split(p[0], terraformPropertyAddrSep)
		}
		return &TerraformPropertyAddr{"", addrs}
	case 2:
		addrs := []string{}
		if p[1] != "" {
			addrs = strings.Split(p[1], terraformPropertyAddrSep)
		}
		return &TerraformPropertyAddr{p[0], addrs}
	default:
		return nil
	}
}

func NewTerraformPropertyAddr(owner, addr string) *TerraformPropertyAddr {
	addrs := []string{}
	if addr != "" {
		addrs = strings.Split(addr, terraformPropertyAddrSep)
	}
	return &TerraformPropertyAddr{owner, addrs}
}

func (addr TerraformPropertyAddr) String() string {
	relative := strings.Join(addr.PropertyAddr, terraformPropertyAddrSep)
	if addr.ResourceName == "" {
		return relative
	}
	return addr.ResourceName + terraformPropertyResourceSep + relative
}

func (addr TerraformPropertyAddr) MarshalJSON() ([]byte, error) {
	return json.Marshal(addr.String())
}

func (addr *TerraformPropertyAddr) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*addr = *ParseTerraformPropertyAddr(s)
	return nil
}

func (addr TerraformPropertyAddr) Append(oaddr string) TerraformPropertyAddr {
	if oaddr == "" {
		return addr
	}
	addrs := make([]string, len(addr.PropertyAddr)+1)
	addrs[len(addrs)-1] = oaddr
	copy(addrs, addr.PropertyAddr)
	return TerraformPropertyAddr{ResourceName: addr.ResourceName, PropertyAddr: addrs}
}

func (addr TerraformPropertyAddr) Contains(oaddr TerraformPropertyAddr) bool {
	if addr.ResourceName != oaddr.ResourceName {
		return false
	}

	if len(addr.PropertyAddr) >= len(oaddr.PropertyAddr) {
		return false
	}
	for i := range addr.PropertyAddr {
		if addr.PropertyAddr[i] != oaddr.PropertyAddr[i] {
			return false
		}
	}
	return true
}

func (addr TerraformPropertyAddr) Equals(oaddr TerraformPropertyAddr) bool {
	if len(addr.PropertyAddr) != len(oaddr.PropertyAddr) || addr.ResourceName != oaddr.ResourceName {
		return false
	}
	for i := range addr.PropertyAddr {
		if addr.PropertyAddr[i] != oaddr.PropertyAddr[i] {
			return false
		}
	}
	return true
}
