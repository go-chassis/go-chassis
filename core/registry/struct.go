package registry

// MicroService struct having full info about micro-service
type MicroService struct {
	ServiceID   string
	AppID       string
	ServiceName string
	Version     string
	Environment string
	Status      string
	Level       string
	Schemas     []string
	Metadata    map[string]string
	Framework   *Framework
	RegisterBy  string
}

// Framework struct having info about micro-service version, name
type Framework struct {
	Name    string
	Version string
}

// MicroServiceInstance struct having full info about micro-service instance
type MicroServiceInstance struct {
	InstanceID      string
	HostName        string
	ServiceID       string
	DefaultProtocol string
	DefaultEndpoint string
	Status          string
	Environment     string
	EndpointsMap    map[string]string
	Metadata        map[string]string
	DataCenterInfo  *DataCenterInfo
	Stage           string
}

// MicroServiceDependency is for to represent dependencies of micro-service
type MicroServiceDependency struct {
	Consumer  *MicroService
	Providers []*MicroService
}

// DataCenterInfo represents micro-service data center info
type DataCenterInfo struct {
	Name          string
	Region        string
	AvailableZone string
}

// SourceInfo source info
type SourceInfo struct {
	Name string
	Tags map[string]string
}

// Schema to represents schema info
type Schema struct {
	Schema string `json:"schema"`
}

// SchemaContent represents schema contents info
type SchemaContent struct {
	Swagger    string                           `yaml:"swagger"`
	Info       map[string]string                `yaml:"info"`
	BasePath   string                           `yaml:"basePath"`
	Produces   []string                         `yaml:"produces"`
	Paths      map[string]map[string]MethodInfo `yaml:"paths"`
	Definition map[string]Definition            `yaml:"definitions"`
}

// SchemaContents represents array of schema contents
type SchemaContents struct {
	Schemas []*SchemaContent
}

// MethodInfo represents method info
type MethodInfo struct {
	OperationID string              `yaml:"operationId"`
	Parameters  []Parameter         `yaml:"parameters"`
	Response    map[string]Response `yaml:"responses"`
}

// Parameter represents schema parameters
type Parameter struct {
	Name      string      `yaml:"name"`
	In        string      `yaml:"in"`
	Required  bool        `yaml:"required"`
	Type      string      `yaml:"type"`
	Format    string      `yaml:"format"`
	Items     Item        `yaml:"items"`
	ColFormat string      `yaml:"collectionFormat"`
	Schema    SchemaValue `yaml:"schema"`
}

// SchemaValue represents additional info of schema
type SchemaValue struct {
	Reference            string                 `yaml:"$ref"`
	Format               string                 `yaml:"format"`
	Title                string                 `yaml:"title"`
	Description          string                 `yaml:"description"`
	Default              string                 `yaml:"default"`
	MultipleOf           int                    `yaml:"multipleOf"`
	ExclusiveMaximum     int                    `yaml:"exclusiveMaximum"`
	Minimum              int                    `yaml:"minimum"`
	ExclusiveMinimum     int                    `yaml:"exclusiveMinimum"`
	MaxLength            int                    `yaml:"maxLength"`
	MinLength            int                    `yaml:"minLength"`
	Pattern              int                    `yaml:"pattern"`
	MaxItems             int                    `yaml:"maxItems"`
	MinItems             int                    `yaml:"minItems"`
	UniqueItems          bool                   `yaml:"uniqueItems"`
	MaxProperties        int                    `yaml:"maxProperties"`
	MinProperties        int                    `yaml:"minProperties"`
	Required             bool                   `yaml:"required"`
	Enum                 []interface{}          `yaml:"enum"`
	Type                 string                 `yaml:"type"`
	Items                Item                   `yaml:"items"`
	Properties           map[string]interface{} `yaml:"properties"`
	AdditionalProperties map[string]string      `yaml:"additionalProperties"`
}

// Item represents type of the schema
type Item struct {
	Type string                 `yaml:"type"`
	XML  map[string]interface{} `yaml:"xml"`
}

// Response represents schema response
type Response struct {
	Description string            `yaml:"description"`
	Schema      map[string]string `yaml:"schema"`
}

// Definition struct represents types, xjavaclass, properities
type Definition struct {
	Types      string                 `yaml:"type"`
	XJavaClass string                 `yaml:"x-java-class"`
	Properties map[string]interface{} `yaml:"properties"`
}
