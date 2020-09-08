package pkg

import (
	"encoding/json"
	"strings"

	"github.com/go-openapi/spec"
)

type propertyAddr struct {
	owner string
	addrs []string
}

const ownerSep = ":"
const addrSep = "."

func newPropertyAddrFromString(addr string) *propertyAddr {
	p := strings.Split(ownerSep, addr)
	switch len(p) {
	case 1:
		return &propertyAddr{"", strings.Split(p[0], addrSep)}
	case 2:
		return &propertyAddr{p[0], strings.Split(p[1], addrSep)}
	default:
		return nil
	}
}

func newPropertyAddr(owner string, addrs ...string) *propertyAddr {
	return &propertyAddr{owner, addrs}
}

func (addr propertyAddr) IsCanonical() bool {
	return addr.owner != ""
}

func (addr propertyAddr) String() string {
	relative := strings.Join(addr.addrs, addrSep)
	if !addr.IsCanonical() {
		return relative
	}
	return addr.owner + ownerSep + relative
}

func (addr propertyAddr) MarshalJSON() ([]byte, error) {
	o := map[string]interface{}{}
	if addr.owner != "" {
		o["owner"] = addr.owner
	}
	o["addr"] = strings.Join(addr.addrs, addrSep)
	return json.Marshal(o)
}

func (addr *propertyAddr) UnmarshalJSON(b []byte) error {
	o := map[string]interface{}{}
	if err := json.Unmarshal(b, &o); err != nil {
		return err
	}
	addr.addrs = strings.Split(o["addr"].(string), addrSep)
	if owner, ok := o["owner"]; ok {
		addr.owner = owner.(string)
	}
	return nil
}

func (addr propertyAddr) Append(oaddr string) propertyAddr {
	addrs := make([]string, len(addr.addrs)+1)
	addrs[len(addrs)-1] = oaddr
	copy(addrs, addr.addrs)
	return propertyAddr{owner: addr.owner, addrs: addrs}
}

func (addr propertyAddr) ToDefinitionRef() (spec.Ref, error) {
	return spec.NewRef("#definitions/" + strings.Join(addr.addrs, "/"))
}
