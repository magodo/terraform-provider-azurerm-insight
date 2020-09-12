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
	return json.Marshal(addr.String())
}

func (addr *PropertyAddr) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*addr = *NewPropertyAddrFromString(s)
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

func (addr PropertyAddr) Equals(oaddr PropertyAddr) bool {
	if len(addr.addrs) != len(oaddr.addrs) || addr.owner != oaddr.owner {
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
	return spec.NewRef(strings.TrimRight("#definitions/"+strings.ReplaceAll(strings.ReplaceAll(addr.String(), ownerSep, "/"), addrSep, "/"), "/"))
}
