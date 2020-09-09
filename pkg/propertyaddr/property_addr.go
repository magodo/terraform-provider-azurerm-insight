package propertyaddr

import (
	"encoding/json"
	"strings"

	"github.com/go-openapi/spec"
)

type PropertyAddr struct {
	owner string // tf: resource name; swg: schema name
	addrs RelativeAddrs
}

type RelativeAddrs []string

func (a RelativeAddrs) String() string {
	return strings.Join(a, addrSep)
}

const ownerSep = ":"
const addrSep = "."

func NewPropertyAddrFromString(addr string) *PropertyAddr {
	p := strings.Split(addr, ownerSep)
	switch len(p) {
	case 1:
		addrs := []string{}
		if p[0] != "" {
			addrs = strings.Split(p[0], addrSep)
		}
		return &PropertyAddr{"", addrs}
	case 2:
		addrs := []string{}
		if p[1] != "" {
			addrs = strings.Split(p[1], addrSep)
		}
		return &PropertyAddr{p[0], addrs}
	default:
		return nil
	}
}

func NewPropertyAddrFromStringWithOwner(owner, addr string) *PropertyAddr {
	addrs := []string{}
	if addr != "" {
		addrs = strings.Split(addr, addrSep)
	}
	return &PropertyAddr{owner, addrs}
}

func NewPropertyAddr(owner string, addrs ...string) *PropertyAddr {
	return &PropertyAddr{owner, addrs}
}

func (addr PropertyAddr) String() string {
	relative := strings.Join(addr.addrs, addrSep)
	if addr.owner == "" {
		return relative
	}
	return addr.owner + ownerSep + relative
}

func (addr PropertyAddr) Owner() string {
	return addr.owner
}

func (addr PropertyAddr) RelativeAddrs() RelativeAddrs {
	return addr.addrs
}

func (addr PropertyAddr) MarshalJSON() ([]byte, error) {
	o := map[string]interface{}{}
	if addr.owner != "" {
		o["owner"] = addr.owner
	}
	if len(addr.addrs) == 0 {
		return json.Marshal(o)
	}
	o["addr"] = strings.Join(addr.addrs, addrSep)
	return json.Marshal(o)
}

func (addr *PropertyAddr) UnmarshalJSON(b []byte) error {
	o := map[string]interface{}{}
	if err := json.Unmarshal(b, &o); err != nil {
		return err
	}
	if a, ok := o["addr"]; ok {
		if a := a.(string); a != "" {
			addr.addrs = strings.Split(a, addrSep)
		}
	}
	if owner, ok := o["owner"]; ok {
		addr.owner = owner.(string)
	}
	return nil
}

func (addr PropertyAddr) Append(oaddr string) PropertyAddr {
	if oaddr == "" {
		return addr
	}
	addrs := make([]string, len(addr.addrs)+1)
	addrs[len(addrs)-1] = oaddr
	copy(addrs, addr.addrs)
	return PropertyAddr{owner: addr.owner, addrs: addrs}
}

func (addr PropertyAddr) Contains(oaddr PropertyAddr) bool {
	if addr.owner != oaddr.owner {
		return false
	}
	// empty addr never contain or is contained by other addr
	if len(addr.addrs) == 0 || len(oaddr.addrs) == 0{
		return false
	}
	if len(addr.addrs) >= len(oaddr.addrs) {
		return false
	}
	for i := range addr.addrs {
		if addr.addrs[i] != oaddr.addrs[i] {
			return false
		}
	}
	return true
}

func ToSwaggerDefinitionRef(addr PropertyAddr) (spec.Ref, error) {
	return spec.NewRef(strings.TrimRight("#definitions/"+ strings.ReplaceAll(strings.ReplaceAll(addr.String(), ownerSep, "/"), addrSep, "/"), "/"))
}
