package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	openapispec "github.com/go-openapi/spec"
	"github.com/magodo/terraform-provider-azurerm-insight/pkg/core/propertyaddr"
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
	TFLinks TFLinks `json:",omitempty"`

	// Whether this property is granted to be not to implement in Terraform
	IsGranted    bool   `json:",omitempty"`
	GrantComment string `json:",omitempty"`

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
	SwaggerRelPath string
	Name           string
	Properties     SWGSchemaProperties

	// Whether this property is granted to be not to implement in Terraform
	IsGranted    bool   `json:",omitempty"`
	GrantComment string `json:",omitempty"`

	swaggerURL    string
	swagger       *openapispec.Swagger
	coverageStore SWGPropertyCoverageStore
}

func NewSWGSchema(swaggerBaseURL, swaggerRelPath string, schemaName string) (*SWGSchema, error) {
	swaggerURI := swaggerBaseURL + "/" + swaggerRelPath
	swagger, err := LoadSwagger(swaggerURI)
	if err != nil {
		return nil, err
	}

	swgSchema := &SWGSchema{
		SwaggerRelPath: swaggerRelPath,
		Name:           schemaName,
		Properties: map[string]*SWGSchemaProperty{
			"": {
				TFLinks: []TFLink{},
				schema:  swagger.Definitions[schemaName],
				resolvedRefs: map[string]interface{}{
					// Consider this schemas itself as resolved reference
					normalizePaths("#/definitions/"+schemaName, swaggerURI): struct{}{},
				},
			},
		},
		swaggerURL: swaggerURI,
		swagger:    swagger,
	}

	// Expand the root level properties of the schemas
	err = swgSchema.ExpandPropertyOneLevelDeep(*propertyaddr.NewPropertyAddrFromStringWithOwner(schemaName, ""))
	if err != nil {
		return nil, fmt.Errorf("expanding schemas %s (%s): %w", schemaName, swaggerURI, err)
	}
	return swgSchema, nil
}

type SWGSchemaCollector func(swagger *openapispec.Swagger) (schemaNames []string)

// CollectSWGSchemas collects the schemas from a swagger spec
func CollectSWGSchemas(swaggerBaseURL, swaggerRelPath string, collector SWGSchemaCollector) ([]SWGSchema, error) {
	swaggerURI := swaggerBaseURL + "/" + swaggerRelPath
	swagger, err := LoadSwagger(swaggerURI)
	if err != nil {
		return nil, err
	}
	schemas := collector(swagger)

	out := make([]SWGSchema, 0, len(schemas))
	for _, schemaName := range schemas {
		schema, err := NewSWGSchema(swaggerBaseURL, swaggerRelPath, schemaName)
		if err != nil {
			return nil, err
		}
		out = append(out, *schema)
	}
	return out, nil
}

// ExpandPropertyOneLevelDeep expand the specified swagger schemas property one level deep, with any allOf and $ref taken into consideration.
func (s *SWGSchema) ExpandPropertyOneLevelDeep(addr propertyaddr.PropertyAddr) error {
	raddr := addr.RelativeAddrs().String()
	prop, ok := s.Properties[raddr]
	if !ok {
		return fmt.Errorf("property %s does not exist in SWGSchema %s (%s)", addr, s.Name, s.swaggerURL)
	}

	isCyclic, err := s.expandRefProperty(prop)
	if err != nil {
		return fmt.Errorf("dereferencing property %s in SWGSchema %s (%s): %w", addr, s.Name, s.swaggerURL, err)
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

		// AllOf contains refs, then need to expandRefProperty then first.

		// We construct a temp SWGSchemaProperty here (as it has no object/property related) to expand it into a concrete schemas.
		// Then we will iterate that schemas's property which by concept is the top level property of this parent property.
		tmpSwgProp := NewSWGSchemaProperty(schema, prop.TFLinks, prop.resolvedRefs)

		isCyclic, err := s.expandRefProperty(tmpSwgProp)
		if err != nil {
			return fmt.Errorf("dereferencing property %s in SWGSchema %s (%s): %w", addr, s.Name, s.swaggerURL, err)
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

// expandRefProperty expand a property itself IN-PLACE until either it is a concrete schemas (i.e. not a ref) or hit a cyclic ref.
func (s *SWGSchema) expandRefProperty(prop *SWGSchemaProperty) (isCyclic bool, err error) {
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

	normalizedRefURI := normalizeFileRef(&ref, s.swaggerURL).String()

	// If current ref has already been derefed, meaning a cyclic ref is hit, we will return.
	if _, ok := prop.resolvedRefs[normalizedRefURI]; ok {
		return true, nil
	}

	// Keep track of the resolved reference to avoid cyclic ref
	prop.resolvedRefs[normalizedRefURI] = struct{}{}

	schema, err := openapispec.ResolveRefWithBase(s.swagger, &ref, &openapispec.ExpandOptions{RelativeBase: s.swaggerURL})
	if err != nil {
		return false, fmt.Errorf("resolve reference %s: %w", ref.String(), err)
	}

	// update the stored schemas by the derefed schemas
	prop.schema = *schema

	return s.expandRefProperty(prop)
}

func (s *SWGSchema) AddTFLink(swgPropAddr, tfPropAddr propertyaddr.PropertyAddr) error {
	var isExpandToChildProperties bool
	for raddr, prop := range s.Properties {
		addr := propertyaddr.NewPropertyAddrFromStringWithOwner(s.Name, raddr)

		if swgPropAddr.Contains(*addr) {
			isExpandToChildProperties = true
			prop.TFLinks = append(prop.TFLinks, TFLink{Prop: tfPropAddr})
			continue
		}

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
	if isExpandToChildProperties {
		return nil
	}
	return fmt.Errorf("property %s doesn't belong to schemas %s (%s)", swgPropAddr, s.Name, s.swaggerURL)
}

// CalcCoverage calculates the property coverage (<=1) of the schema/property, and fill in the SWGSchema.
// Those granted properties are not counted during the calculation.
func (s *SWGSchema) CalcCoverage() error {
	store := NewSWGPropertyCoverageStore()
	for propAddr, prop := range s.Properties {
		if err := store.Add(*propertyaddr.NewPropertyAddrFromString(propAddr), *prop); err != nil {
			return fmt.Errorf("adding property %q: %v", propAddr, err)
		}
	}
	s.coverageStore = store
	return nil
}

func (s *SWGSchema) SchemaCoverage() (covered, total int) {
	return s.coverageStore.SchemaCoverage()
}

func (s *SWGSchema) FindCoverage(propAddr propertyaddr.PropertyAddr) (covered, total int, ok bool) {
	return s.coverageStore.FindCoverage(propAddr)
}

const swgSchemaAddrSep = "#/definitions/"

type SWGSchemaAddr string

func NewSWGSchemaAddr(swaggerRelPath, schemaName string) SWGSchemaAddr {
	return SWGSchemaAddr(swaggerRelPath + swgSchemaAddrSep + schemaName)
}

func (addr SWGSchemaAddr) SwaggerRelPath() string {
	return strings.Split(string(addr), swgSchemaAddrSep)[0]
}

func (addr SWGSchemaAddr) SchemaName() string {
	return strings.Split(string(addr), swgSchemaAddrSep)[1]
}

// SWGSchemas caches the SWGSchema using swagger + schemas as key.
// During each link operation from terraform schemas to swagger schemas, it will manipulate one of
// the SWGSchema. Afterwards, this type contains all the mapping info from swagger to terraform.
type SWGSchemas struct {
	sync.Mutex
	m map[SWGSchemaAddr]*SWGSchema
}

func (c *SWGSchemas) Lock() {
	c.Mutex.Lock()
}

func (c *SWGSchemas) Unlock() {
	c.Mutex.Unlock()
}

func (c *SWGSchemas) Get(addr SWGSchemaAddr) *SWGSchema {
	return c.m[addr]
}

func (c *SWGSchemas) Set(addr SWGSchemaAddr, schema *SWGSchema) {
	c.m[addr] = schema
}

func NewSGWSchemas() *SWGSchemas {
	return &SWGSchemas{
		Mutex: sync.Mutex{},
		m:     map[SWGSchemaAddr]*SWGSchema{},
	}
}

// Build SWGSchemas by processing on Terraform schema files (which resides in the tfSchemaDir)
// and the Swagger specs (which resides in the swaggerBaseDir, can be either a local path or an http URI)
// Optionally, users can specify the swaggerGrantDir which contains the grants for those non-terraform
// appropriate swagger schema/properties.
func NewSWGSchemasFromTerraformSchema(swaggerBasePath, tfSchemaDir, swaggerGrantBaseDir string) (*SWGSchemas, error) {
	swgschemas := NewSGWSchemas()
	err := filepath.Walk(tfSchemaDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		var tfschema TFSchema
		if err := json.Unmarshal(b, &tfschema); err != nil {
			return err
		}
		if err := tfschema.Validate(); err != nil {
			return fmt.Errorf("validating tf schema %s: %v", tfschema.Name, err)
		}

		if err := tfschema.LinkSwagger(*swgschemas, swaggerBasePath); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking the terraform schema directory %q: %v\n", tfSchemaDir, err)
	}

	// grant swagger schemas
	if swaggerGrantBaseDir != "" {
		swggrant, err := NewSWGGrantFromFiles(swaggerGrantBaseDir)
		if err != nil {
			log.Fatal(err)
		}
		if err := swgschemas.Grant(swggrant); err != nil {
			return nil, err
		}
	}

	// calculate swagger property coverage
	for schemaAddr, schema := range swgschemas.GetAll() {
		if err := schema.CalcCoverage(); err != nil {
			log.Fatalf("calculating coverage for %q: %v", schemaAddr, err)
		}
	}
	return swgschemas, nil
}

func (c *SWGSchemas) LinkSWGSchema(swaggerBasePath, swaggerRelPath string, swgPropAddr, tfPropAddr propertyaddr.PropertyAddr) error {
	c.Lock()
	defer c.Unlock()

	swgSchema := c.Get(NewSWGSchemaAddr(swaggerRelPath, swgPropAddr.Owner()))
	if swgSchema == nil {
		var err error
		swgSchema, err = NewSWGSchema(swaggerBasePath, swaggerRelPath, swgPropAddr.Owner())
		if err != nil {
			return err
		}
	}

	defer c.Set(NewSWGSchemaAddr(swaggerRelPath, swgPropAddr.Owner()), swgSchema)

	return swgSchema.AddTFLink(swgPropAddr, tfPropAddr)
}

// Grant inquiries the SWGGrant to add the granting information onto the SWGSchemas
func (c *SWGSchemas) Grant(grant SWGGrant) error {
	c.Lock()
	defer c.Unlock()
	for schemaAddr, schemaGrant := range grant {
		schema, ok := c.m[schemaAddr]
		if !ok {
			continue
		}

		if schemaGrant.IsSchemaGranted() {
			schema.IsGranted = true
			schema.GrantComment = schemaGrant.Comment
			continue
		}

		for propertyAddr, propertyGrantComment := range schemaGrant.Properties {
			property, ok := schema.Properties[propertyAddr]
			if !ok {
				return fmt.Errorf(`property to be granted: "%s" doesn't exist in Swagger schema: %s'`, propertyAddr, schemaAddr)
			}
			property.IsGranted = true
			property.GrantComment = propertyGrantComment
		}
	}
	return nil
}

// GetSWGSchema get all SWGSchema from cache.
func (c *SWGSchemas) GetAll() map[SWGSchemaAddr]*SWGSchema {
	c.Lock()
	defer c.Unlock()
	out := map[SWGSchemaAddr]*SWGSchema{}
	for k, v := range c.m {
		out[k] = v
	}
	return out
}
