package pkg

import (
	"encoding/json"
	"fmt"
	"sync"

	openapispec "github.com/go-openapi/spec"
	"github.com/magodo/terraform-provider-azurerm-insight/pkg/propertyaddr"
)

type TFLink struct {
	Prop propertyaddr.PropertyAddr
}

type TFLinks []TFLink

func (links TFLinks) MarshalJSON() ([]byte, error) {
	addrs := []string{}
	for _, link := range links {
		addrs = append(addrs, link.Prop.String())
	}
	return json.Marshal(addrs)
}

func (links *TFLinks) UnmarshalJSON(b []byte) error {
	var addrs []string
	if err := json.Unmarshal(b, &addrs); err != nil {
		return err
	}
	*links = []TFLink{}
	for _, addr := range addrs {
		*links = append(*links, TFLink{*propertyaddr.NewPropertyAddrFromString(addr)})
	}
	return nil
}

type SWGSchemaProperty struct {
	// Terraform property addresses
	TFLinks TFLinks

	// The schemas of this swagger schemas property
	schema openapispec.Schema

	// The resolved URI refs along the way to this schemas, each is an absolute/normalized reference.
	resolvedRefs map[string]interface{}
}

func NewSWGSchemaProperty(schema openapispec.Schema, tflinks []TFLink, resolvedRefs map[string]interface{}) *SWGSchemaProperty {
	newTFLinks := []TFLink{}
	if len(tflinks) != 0 {
		newTFLinks = make([]TFLink, len(tflinks))
		copy(newTFLinks, tflinks)
	}

	newResolvedRefs := map[string]interface{}{}
	if len(resolvedRefs) != 0 {
		for k, v := range resolvedRefs {
			newResolvedRefs[k] = v
		}
	}
	return &SWGSchemaProperty{
		TFLinks:      newTFLinks,
		schema:       schema,
		resolvedRefs: newResolvedRefs,
	}
}

type SWGSchemaProperties map[string]*SWGSchemaProperty // the key is swagger schemas relative property addr

type SWGSchema struct {
	Name       string
	SpecPath   string
	Properties SWGSchemaProperties

	swagger *openapispec.Swagger
}

func NewSWGSchema(specPath string, schemaName string) (*SWGSchema, error) {
	swagger, err := LoadSwagger(specPath)
	if err != nil {
		return nil, err
	}

	swgSchema := &SWGSchema{
		Name:     schemaName,
		SpecPath: specPath,
		Properties: map[string]*SWGSchemaProperty{
			"": {
				TFLinks: []TFLink{},
				schema:  swagger.Definitions[schemaName],
				resolvedRefs: map[string]interface{}{
					// Consider this schemas itself as resolved reference
					normalizePaths("#/definitions/"+schemaName, specPath): struct{}{},
				},
			},
		},
		swagger: swagger,
	}

	// Expand the root level properties of the schemas
	err = swgSchema.ExpandPropertyOneLevelDeep(*propertyaddr.NewPropertyAddrFromStringWithOwner(schemaName, ""))
	if err != nil {
		return nil, fmt.Errorf("expanding schemas %s (%s): %w", schemaName, specPath, err)
	}
	return swgSchema, nil
}

// ExpandPropertyOneLevelDeep expand the specified swagger schemas property one level deep, with any allOf and $ref taken into consideration.
func (s *SWGSchema) ExpandPropertyOneLevelDeep(addr propertyaddr.PropertyAddr) error {
	raddr := addr.RelativeAddrs().String()
	prop, ok := s.Properties[raddr]
	if !ok {
		return fmt.Errorf("property %s does not exist in SWGSchema %s (%s)", addr, s.Name, s.SpecPath)
	}

	isCyclic, err := s.expandProperty(prop)
	if err != nil {
		return fmt.Errorf("dereferencing property %s in SWGSchema %s (%s): %w", addr, s.Name, s.SpecPath, err)
	}

	// If the property to be expanded is a cyclic reference, we will do nothing but keep that property
	if isCyclic {
		return nil
	}

	// direct top level properties
	for propK, propV := range prop.schema.Properties {
		p := NewSWGSchemaProperty(propV, prop.TFLinks, prop.resolvedRefs)
		addr := addr.Append(propK)
		s.addProperty(addr, *p)
	}

	// expand AllOf properties
	for _, schema := range prop.schema.AllOf {

		// AllOf contains concrete schemas, then directly add the property.
		if schema.Ref.String() == "" {
			for propK, propV := range schema.Properties {
				p := NewSWGSchemaProperty(propV, prop.TFLinks, prop.resolvedRefs)
				addr := addr.Append(propK)
				s.addProperty(addr, *p)
			}
			continue
		}

		// AllOf contains refs, then need to expandProperty then first.

		// We construct a temp SWGSchemaProperty here (as it has no object/property related) to expand it into a concrete schemas.
		// Then we will iterate that schemas's property which by concept is the top level property of this parent property.
		tmpSwgProp := NewSWGSchemaProperty(schema, prop.TFLinks, prop.resolvedRefs)

		isCyclic, err := s.expandProperty(tmpSwgProp)
		if err != nil {
			return fmt.Errorf("dereferencing property %s in SWGSchema %s (%s): %w", addr, s.Name, s.SpecPath, err)
		}

		// Ignore as there is no better way to handle this (since it has no object/property related)
		if isCyclic {
			continue
		}

		for propK, propV := range tmpSwgProp.schema.Properties {
			p := NewSWGSchemaProperty(propV, tmpSwgProp.TFLinks, tmpSwgProp.resolvedRefs)
			addr := addr.Append(propK)
			s.addProperty(addr, *p)
		}
	}

	// We have to check whether we added any child property of this property. If this property is already the leaf property,
	// we should keep this property from removing it from the SWGSchema property map.
	for currentRAddr := range s.Properties {
		currentAddr := propertyaddr.NewPropertyAddrFromStringWithOwner(s.Name, currentRAddr)
		if addr.Contains(*currentAddr) {
			delete(s.Properties, raddr)
			return nil
		}
	}
	return nil
}

// addProperty adds a new SWGSchemaProperty to the SWGSchema.
func (s *SWGSchema) addProperty(addr propertyaddr.PropertyAddr, prop SWGSchemaProperty) {
	s.Properties[addr.RelativeAddrs().String()] = &prop
}

// expandProperty expand a property itself IN-PLACE until either it is a concrete schemas (i.e. not a ref) or hit a cyclic ref.
func (s *SWGSchema) expandProperty(prop *SWGSchemaProperty) (isCyclic bool, err error) {
	ref := prop.schema.Ref
	if ref.String() == "" {
		// Specially, if current schema is an array and the items is a ref, we need to go on expand it.
		if prop.schema.Items == nil {
			return false, nil
		}

		if prop.schema.Items.Schema == nil || len(prop.schema.Items.Schemas) != 0 {
			return false, nil
		}

		var schema openapispec.Schema
		if prop.schema.Items.Schema != nil {
			schema = *prop.schema.Items.Schema
		} else {
			schema = prop.schema.Items.Schemas[0]
		}
		if schema.Ref.String() == "" {
			return false, nil
		}
		// continue expanding the ref of the array item
		ref = schema.Ref
	}

	normalizedRefURI := normalizeFileRef(&ref, s.SpecPath).String()

	// If current ref has already been derefed, meaning a cyclic ref is hit, we will return.
	if _, ok := prop.resolvedRefs[normalizedRefURI]; ok {
		return true, nil
	}

	// Keep track of the resolved reference to avoid cyclic ref
	prop.resolvedRefs[normalizedRefURI] = struct{}{}

	schema, err := openapispec.ResolveRefWithBase(s.swagger, &ref, &openapispec.ExpandOptions{RelativeBase: s.SpecPath})
	if err != nil {
		return false, fmt.Errorf("resolve reference %s: %w", ref.String(), err)
	}

	// update the stored schemas by the derefed schemas
	prop.schema = *schema

	return s.expandProperty(prop)
}

func (s *SWGSchema) AddTFLink(swgPropAddr, tfPropAddr propertyaddr.PropertyAddr) error {
	for raddr, prop := range s.Properties {
		addr := propertyaddr.NewPropertyAddrFromStringWithOwner(s.Name, raddr)

		if !addr.Contains(swgPropAddr) && !addr.Equals(swgPropAddr) {
			continue
		}

		if addr.Equals(swgPropAddr) {
			prop.TFLinks = append(prop.TFLinks, TFLink{Prop: tfPropAddr})
			return nil
		}

		// The schemas property we're seeking is a direct or indirect member of the property under iteration
		if err := s.ExpandPropertyOneLevelDeep(*addr); err != nil {
			return fmt.Errorf("expanding top level property for %s: %w", addr, err)
		}
		return s.AddTFLink(swgPropAddr, tfPropAddr)
	}
	return fmt.Errorf("property %s doesn't belong to schemas %s (%s)", swgPropAddr, s.Name, s.SpecPath)
}

type SWGSpecSchemaCache struct {
	sync.Mutex
	m map[string]*SWGSchema
}

func (c *SWGSpecSchemaCache) Lock() {
	c.Mutex.Lock()
}

func (c *SWGSpecSchemaCache) Unlock() {
	c.Mutex.Unlock()
}

func (c *SWGSpecSchemaCache) Get(specPath, schemaName string) *SWGSchema {
	k := specPath + "#/definitions/" + schemaName
	return c.m[k]
}

func (c *SWGSpecSchemaCache) Set(specPath, schemaName string, schema *SWGSchema) {
	k := specPath + "#/definitions/" + schemaName
	c.m[k] = schema
}

// swgSpecSchemaCache caches the SWGSchema using swagger + schemas as key.
// During each link operation from terraform schemas to swagger schemas, it will manipulate one of
// the SWGSchema in this cache. Afterwards, this cache contains all the mapping info from swagger to terraform.
var swgSpecSchemaCache = SWGSpecSchemaCache{
	Mutex: sync.Mutex{},
	m:     map[string]*SWGSchema{},
}

func LinkSWGSchema(specPath string, swgPropAddr, tfPropAddr propertyaddr.PropertyAddr) error {
	swgSpecSchemaCache.Lock()
	defer swgSpecSchemaCache.Unlock()

	swgSchema := swgSpecSchemaCache.Get(specPath, swgPropAddr.Owner())
	if swgSchema == nil {
		var err error
		swgSchema, err = NewSWGSchema(specPath, swgPropAddr.Owner())
		if err != nil {
			return err
		}
	}

	defer swgSpecSchemaCache.Set(specPath, swgPropAddr.Owner(), swgSchema)

	return swgSchema.AddTFLink(swgPropAddr, tfPropAddr)
}

// GetSWGSchema get all SWGSchema from cache.
func GetAllSWGSchemas() map[string]*SWGSchema {
	swgSpecSchemaCache.Lock()
	defer swgSpecSchemaCache.Unlock()
	out := map[string]*SWGSchema{}
	for k, v := range swgSpecSchemaCache.m {
		out[k] = v
	}
	return out
}
