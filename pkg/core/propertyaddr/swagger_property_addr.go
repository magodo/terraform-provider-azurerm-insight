package propertyaddr

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/magodo/terraform-provider-azurerm-insight/pkg/core/utils"

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

func MustParseSwaggerPropertyAddr(addr string) SwaggerPropertyAddr {
	swgPropAddr, err := ParseSwaggerPropertyAddr(addr)
	if err != nil {
		panic(err)
	}

	return swgPropAddr
}

func ParseSwaggerPropertyAddr(addr string) (SwaggerPropertyAddr, error) {
	if strings.Trim(addr, " ") == "" {
		return SwaggerPropertyAddr{}, nil
	}

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

	props, err := ParseSwaggerRelPropertyAddr(propAddr)
	if err != nil {
		return SwaggerPropertyAddr{}, err
	}

	return SwaggerPropertyAddr{
		Schema:       schemaName,
		PropertyAddr: props,
	}, nil
}

func MustNewSwaggerPropertyAddr(schemaName string, propAddr string) SwaggerPropertyAddr {
	addr, err := NewSwaggerPropertyAddr(schemaName, propAddr)
	if err != nil {
		panic(err)
	}
	return addr
}

func NewSwaggerPropertyAddr(schemaName string, propAddr string) (SwaggerPropertyAddr, error) {
	props, err := ParseSwaggerRelPropertyAddr(propAddr)
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

	if len(addr.PropertyAddr) > len(oaddr.PropertyAddr) {
		return false
	}

	if len(addr.PropertyAddr) == len(oaddr.PropertyAddr) {
		lastIdx := len(addr.PropertyAddr) - 1
		if lastIdx == -1 {
			return false
		}
		addrLastProp := addr.PropertyAddr[lastIdx]
		oaddrLastProp := oaddr.PropertyAddr[lastIdx]
		if addrLastProp.discriminatorValue != nil || oaddrLastProp.discriminatorValue == nil {
			return false
		}
	}

	for i := range addr.PropertyAddr {
		addrProp := addr.PropertyAddr[i]
		oaddrProp := oaddr.PropertyAddr[i]
		if addrProp.name != oaddrProp.name {
			return false
		}

		if addrProp.discriminatorValue != nil && oaddrProp.discriminatorValue == nil ||
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
	if addr.Schema == "" {
		return spec.Ref{}, fmt.Errorf("can't turn into swagger definition reference since the schema name is empty: %s", addr.String())
	}

	ref := "#definitions/" + addr.Schema

	for _, prop := range addr.PropertyAddr {
		if prop.discriminatorValue != nil {
			ref += "/" + *prop.discriminatorValue
			continue
		}
		ref += "/" + prop.name
	}
	return spec.NewRef(ref)
}

func (addr SwaggerPropertyAddr) Append(oaddr string) (SwaggerPropertyAddr, error) {
	if oaddr == "" {
		return addr, nil
	}

	oPropAddrs, err := ParseSwaggerRelPropertyAddr(oaddr)
	if err != nil {
		return SwaggerPropertyAddr{}, err
	}

	newPropAddrs := make(SwaggerRelPropertyAddr, len(addr.PropertyAddr)+len(oPropAddrs))
	copy(newPropAddrs, addr.PropertyAddr)
	copy(newPropAddrs[len(addr.PropertyAddr):], oPropAddrs)
	return SwaggerPropertyAddr{
		Schema:       addr.Schema,
		PropertyAddr: newPropAddrs,
	}, nil
}

func (addr SwaggerPropertyAddr) SetDiscriminator(variant string) {
	addr.PropertyAddr[len(addr.PropertyAddr)-1].discriminatorValue = &variant
	return
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

type SwaggerRelPropertyAddr []SwaggerPropertyAddrSegment

type SwaggerPropertyAddrSegment struct {
	name string

	// discriminatorValue is non-nil only when the property is a derived model, and is nil otherwise.
	discriminatorValue *string
}

func NewSwaggerPropertyAddrSegment(name string, discriminator *string) SwaggerPropertyAddrSegment {
	return SwaggerPropertyAddrSegment{
		name:               name,
		discriminatorValue: discriminator,
	}
}

func (prop SwaggerPropertyAddrSegment) String() string {
	v := prop.name
	if prop.discriminatorValue != nil {
		v += swaggerPropertyDiscriminatorStartMark + *prop.discriminatorValue + swaggerPropertyDiscriminatorEndMark
	}
	return v
}

func ParseSwaggerRelPropertyAddr(addr string) (SwaggerRelPropertyAddr, error) {
	if strings.Trim(addr, " ") == "" {
		return nil, nil
	}
	var props SwaggerRelPropertyAddr
	discriminatorPattern := regexp.MustCompile(fmt.Sprintf(`^(.+)?\%s(.+)\%s$`, swaggerPropertyDiscriminatorStartMark, swaggerPropertyDiscriminatorEndMark))
	for _, prop := range strings.Split(addr, swaggerPropertyAddrSep) {
		if !strings.HasSuffix(prop, swaggerPropertyDiscriminatorEndMark) {
			props = append(props, SwaggerPropertyAddrSegment{
				name: prop,
			})
			continue
		}

		m := discriminatorPattern.FindStringSubmatch(prop)
		if len(m) != 3 {
			return nil, fmt.Errorf(`invalid discriminator property notation: %q (expected format: "prop[variant]")`, prop)
		}
		props = append(props, SwaggerPropertyAddrSegment{
			name:               m[1],
			discriminatorValue: &m[2],
		})
	}
	return props, nil
}

func (addr SwaggerRelPropertyAddr) String() string {
	props := []string{}
	for _, prop := range addr {
		props = append(props, prop.String())
	}

	return strings.Join(props, swaggerPropertyAddrSep)
}

func (addr SwaggerPropertyAddr) AsVariant(v string) SwaggerPropertyAddr {
	newAddr := addr.Copy()
	if len(newAddr.PropertyAddr) == 0 {
		newAddr.PropertyAddr = SwaggerRelPropertyAddr{
			{
				discriminatorValue: &v,
			},
		}
		return newAddr
	}
	newAddr.PropertyAddr[len(newAddr.PropertyAddr)-1].discriminatorValue = &v
	return newAddr
}

func (addr SwaggerPropertyAddr) Copy() SwaggerPropertyAddr {
	newAddr := SwaggerPropertyAddr{
		Schema:       addr.Schema,
		PropertyAddr: SwaggerRelPropertyAddr{},
	}

	for _, segment := range addr.PropertyAddr {
		newSegment := SwaggerPropertyAddrSegment{
			name: segment.name,
		}
		if segment.discriminatorValue != nil {
			newSegment.discriminatorValue = utils.String(*segment.discriminatorValue)
		}
		newAddr.PropertyAddr = append(newAddr.PropertyAddr, newSegment)
	}
	return newAddr
}
