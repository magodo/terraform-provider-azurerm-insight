package propertyaddr

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-openapi/spec"
)

const swaggerPropertySchemaSep = ":"
const swaggerPropertyAddrSep = "."
const swaggerPropertyDiscriminatorStartMark = "["
const swaggerPropertyDiscriminatorEndMark = "]"

type SwaggerPropertyAddr struct {
	Schema       string
	PropertyAddr SwaggerRelPropertyAddr
}

type SwaggerRelPropertyAddr []SwaggerProperty

type SwaggerProperty struct {
	name string

	// discriminatorValue is non-nil only when the property is a derived model, and is nil otherwise.
	discriminatorValue *string
}

func ParseSwaggerPropertyAddr(addr string) (SwaggerPropertyAddr, error) {
	p := strings.Split(addr, swaggerPropertySchemaSep)
	var (
		schemaName string
		propAddr   string
	)
	switch len(p) {
	case 1:
		propAddr = p[0]
	case 2:
		schemaName = p[0]
		propAddr = p[1]
	default:
		return SwaggerPropertyAddr{}, fmt.Errorf("invalid Swagger Property Address: %s", addr)
	}

	props, err := parseSwaggerRelPropertyAddr(propAddr)
	if err != nil {
		return SwaggerPropertyAddr{}, err
	}

	return SwaggerPropertyAddr{
		Schema:       schemaName,
		PropertyAddr: props,
	}, nil
}

func parseSwaggerRelPropertyAddr(addr string) (SwaggerRelPropertyAddr, error) {
	var props SwaggerRelPropertyAddr
	discriminatorPattern := regexp.MustCompile(fmt.Sprintf(`^(.+)\%s(.+)\%s$`, swaggerPropertyDiscriminatorStartMark, swaggerPropertyDiscriminatorEndMark))
	for _, prop := range strings.Split(addr, swaggerPropertyAddrSep) {
		if !strings.HasSuffix(prop, swaggerPropertyDiscriminatorEndMark) {
			props = append(props, SwaggerProperty{
				name: prop,
			})
			continue
		}

		m := discriminatorPattern.FindStringSubmatch(prop)
		if len(m) != 3 {
			return nil, fmt.Errorf(`invalid discriminator property notation: %q (expected format: "prop[variant]")`, prop)
		}
		props = append(props, SwaggerProperty{
			name:               m[1],
			discriminatorValue: &m[2],
		})
	}
	return props, nil
}

func NewSwaggerPropertyAddr(schemaName string, propAddr string) (SwaggerPropertyAddr, error) {
	props, err := parseSwaggerRelPropertyAddr(propAddr)
	if err != nil {
		return SwaggerPropertyAddr{}, err
	}
	return SwaggerPropertyAddr{
		Schema:       schemaName,
		PropertyAddr: props,
	}, nil
}

func (addr SwaggerPropertyAddr) Contains(oaddr SwaggerPropertyAddr) bool {
	if addr.Schema != oaddr.Schema {
		return false
	}

	if len(addr.PropertyAddr) >= len(oaddr.PropertyAddr) {
		return false
	}

	for i := range addr.PropertyAddr {
		addrProp := addr.PropertyAddr[i]
		oaddrProp := oaddr.PropertyAddr[i]
		if addrProp.name != oaddrProp.name {
			return false
		}
		if addrProp.discriminatorValue == nil && oaddrProp.discriminatorValue != nil ||
			addrProp.discriminatorValue != nil && oaddrProp.discriminatorValue == nil ||
			addrProp.discriminatorValue != nil && oaddrProp.discriminatorValue != nil && *addrProp.discriminatorValue != *oaddrProp.discriminatorValue {
			return false
		}
	}

	return true
}

func (addr SwaggerPropertyAddr) Equals(oaddr SwaggerPropertyAddr) bool {
	if addr.Schema != oaddr.Schema {
		return false
	}

	if len(addr.PropertyAddr) != len(oaddr.PropertyAddr) {
		return false
	}

	for i := range addr.PropertyAddr {
		addrProp := addr.PropertyAddr[i]
		oaddrProp := oaddr.PropertyAddr[i]
		if addrProp.name != oaddrProp.name {
			return false
		}
		if addrProp.discriminatorValue == nil && oaddrProp.discriminatorValue != nil ||
			addrProp.discriminatorValue != nil && oaddrProp.discriminatorValue == nil ||
			addrProp.discriminatorValue != nil && oaddrProp.discriminatorValue != nil && *addrProp.discriminatorValue != *oaddrProp.discriminatorValue {
			return false
		}
	}

	return true
}

func (addr SwaggerPropertyAddr) ToSwaggerDefinitionRef() (spec.Ref, error) {
	// TODO: can't directly use the addr.String() as it contains special discriminator notation (foo[variant]).
	return spec.NewRef(strings.TrimRight("#definitions/"+strings.ReplaceAll(strings.ReplaceAll(addr.String(), swaggerPropertySchemaSep, "/"), swaggerPropertyAddrSep, "/"), "/"))
}

func (addr SwaggerPropertyAddr) String() string {
	if addr.Schema == "" {
		return addr.PropertyAddr.String()
	}
	return addr.Schema + swaggerPropertySchemaSep + addr.PropertyAddr.String()
}

func (addr SwaggerPropertyAddr) MarshalJSON() ([]byte, error) {
	return json.Marshal(addr.String())
}

func (addr *SwaggerPropertyAddr) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	var err error
	*addr, err = ParseSwaggerPropertyAddr(s)
	return err
}

func (addr SwaggerRelPropertyAddr) String() string {
	props := []string{}
	for _, prop := range addr {
		v := prop.name
		if prop.discriminatorValue != nil {
			v += swaggerPropertyDiscriminatorStartMark + *prop.discriminatorValue + swaggerPropertyDiscriminatorEndMark
		}
		props = append(props, v)
	}

	return strings.Join(props, swaggerPropertyAddrSep)
}

func (addr SwaggerRelPropertyAddr) Append(oaddr string) (SwaggerRelPropertyAddr, error) {
	if oaddr == "" {
		return addr, nil
	}

	oaddrs, err := parseSwaggerRelPropertyAddr(oaddr)
	if err != nil {
		return SwaggerRelPropertyAddr{}, err
	}

	newaddr := make(SwaggerRelPropertyAddr, len(addr)+len(oaddrs))
	copy(newaddr, addr)
	copy(newaddr[len(addr):], oaddrs)
	return newaddr, nil
}
