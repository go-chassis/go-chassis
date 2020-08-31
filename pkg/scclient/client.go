package client

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/apache/servicecomb-service-center/pkg/registry"
	"github.com/cenkalti/backoff"
	"github.com/go-chassis/foundation/httpclient"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
	"github.com/go-chassis/go-chassis/resilience/retry"
	"github.com/go-chassis/openlog"
	"github.com/gorilla/websocket"
)

// Define constants for the client
const (
	MicroservicePath    = "/microservices"
	InstancePath        = "/instances"
	BatchInstancePath   = "/instances/action"
	SchemaPath          = "/schemas"
	HeartbeatPath       = "/heartbeat"
	ExistencePath       = "/existence"
	WatchPath           = "/watcher"
	StatusPath          = "/status"
	DependencyPath      = "/dependencies"
	PropertiesPath      = "/properties"
	HeaderContentType   = "Content-Type"
	HeaderUserAgent     = "User-Agent"
	DefaultAddr         = "127.0.0.1:30100"
	AppsPath            = "/apps"
	DefaultRetryTimeout = 500 * time.Millisecond
	HeaderRevision      = "X-Resource-Revision"
	EnvProjectID        = "CSE_PROJECT_ID"
	// EnvCheckSCIInterval sc instance health check interval in second
	EnvCheckSCIInterval = "CHASSIS_SC_HEALTH_CHECK_INTERVAL"
)

// Define variables for the client
var (
	MSAPIPath     = ""
	GovernAPIPATH = ""
	TenantHeader  = "X-Domain-Name"
)
var (
	//ErrNotModified means instance is not changed
	ErrNotModified = errors.New("instance is not changed since last query")
	//ErrMicroServiceExists means service is registered
	ErrMicroServiceExists = errors.New("micro-service already exists")
	// ErrMicroServiceNotExists means service is not exists
	ErrMicroServiceNotExists = errors.New("micro-service does not exist")
	//ErrEmptyCriteria means you gave an empty list of criteria
	ErrEmptyCriteria = errors.New("batch find criteria is empty")
)

// RegistryClient is a structure for the client to communicate to Service-Center
type RegistryClient struct {
	Config     *RegistryConfig
	client     *httpclient.Requests
	protocol   string
	watchers   map[string]bool
	mutex      sync.Mutex
	wsDialer   *websocket.Dialer
	conns      map[string]*websocket.Conn
	apiVersion string
	revision   string
	pool       *AddressPool
}

// RegistryConfig is a structure to store registry configurations like address of cc, ssl configurations and tenant name
type RegistryConfig struct {
	SSL bool
}

// URLParameter maintains the list of parameters to be added in URL
type URLParameter map[string]string

//ResetRevision reset the revision to 0
func (c *RegistryClient) ResetRevision() {
	c.revision = "0"
}

// Initialize initializes the Registry Client
func (c *RegistryClient) Initialize(opt Options) (err error) {
	c.revision = "0"
	c.Config = &RegistryConfig{
		SSL: opt.EnableSSL,
	}

	options := &httpclient.Options{
		SSLEnabled: opt.EnableSSL,
		TLSConfig:  opt.TLSConfig,
		Compressed: opt.Compressed,
	}
	c.watchers = make(map[string]bool)
	c.conns = make(map[string]*websocket.Conn)
	c.protocol = "https"
	c.wsDialer = &websocket.Dialer{
		TLSClientConfig: opt.TLSConfig,
	}
	if !c.Config.SSL {
		c.wsDialer = websocket.DefaultDialer
		c.protocol = "http"
	}
	c.client, err = httpclient.New(options)
	if err != nil {
		return err
	}

	//Set the API Version based on the value set in chassis.yaml servicecomb.registry.api.version
	//Default Value Set to V4
	opt.Version = strings.ToLower(opt.Version)
	switch opt.Version {
	case "v3":
		c.apiVersion = "v3"
		TenantHeader = "X-Tenant-Name"
	default:
		c.apiVersion = "v4"
	}
	//Update the API Base Path based on the Version
	c.updateAPIPath()
	c.pool = GetInstance()
	c.pool.SetAddress(opt.Addrs)
	return nil
}

// updateAPIPath Updates the Base PATH anf HEADERS Based on the version of SC used.
func (c *RegistryClient) updateAPIPath() {
	//Check for the env Name in Container to get Domain Name
	//Default value is  "default"
	projectID, isExist := os.LookupEnv(EnvProjectID)
	if !isExist {
		projectID = "default"
	}
	switch c.apiVersion {
	case "v3":
		MSAPIPath = APIPath
		GovernAPIPATH = APIPath
		openlog.Info("Use Service center v3")
	default:
		MSAPIPath = "/v4/" + projectID + "/registry"
		GovernAPIPATH = "/v4/" + projectID + "/govern"
		openlog.Info("Use Service center v4")
	}
}

// SyncEndpoints gets the endpoints of service-center in the cluster
//you only need to call this function,
//if your service center is not behind a load balancing service like ELB,nginx etc
func (c *RegistryClient) SyncEndpoints() error {
	c.pool.Monitor()
	instances, err := c.Health()
	if err != nil {
		return fmt.Errorf("sync SC ep failed. err:%s", err.Error())
	}
	eps := make([]string, 0)
	for _, instance := range instances {
		m := getProtocolMap(instance.Endpoints)
		eps = append(eps, m["rest"])
	}
	if len(eps) != 0 {
		c.pool.SetAddress(eps)
		openlog.Info("Sync service center endpoints " + strings.Join(eps, ","))
		return nil
	}
	return fmt.Errorf("sync endpoints failed")
}

func (c *RegistryClient) formatURL(api string, querys []URLParameter, options *CallOptions) string {
	builder := URLBuilder{
		Protocol:      c.protocol,
		Host:          c.getAddress(),
		Path:          api,
		URLParameters: querys,
		CallOptions:   options,
	}
	return builder.String()
}

// GetDefaultHeaders gets the default headers for each request to be made to Service-Center
func (c *RegistryClient) GetDefaultHeaders() http.Header {

	headers := http.Header{
		HeaderContentType: []string{"application/json"},
		HeaderUserAgent:   []string{"cse-serviceregistry-client/1.0.0"},
		TenantHeader:      []string{"default"},
	}

	return headers
}

// HTTPDo makes the http request to Service-center with proper header, body and method
func (c *RegistryClient) HTTPDo(method string, rawURL string, headers http.Header, body []byte) (resp *http.Response, err error) {
	if len(headers) == 0 {
		headers = make(http.Header)
	}
	for k, v := range c.GetDefaultHeaders() {
		headers[k] = v
	}
	return c.client.Do(context.Background(), method, rawURL, headers, body)
}

// RegisterService registers the micro-services to Service-Center
func (c *RegistryClient) RegisterService(microService *registry.MicroService) (string, error) {
	if microService == nil {
		return "", errors.New("invalid request MicroService parameter")
	}
	request := registry.CreateServiceRequest{
		Service: microService,
	}

	registerURL := c.formatURL(MSAPIPath+MicroservicePath, nil, nil)
	body, err := json.Marshal(request)
	if err != nil {
		return "", NewJSONException(err, string(body))
	}

	resp, err := c.HTTPDo("POST", registerURL, nil, body)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", fmt.Errorf("RegisterService failed, response is empty, MicroServiceName: %s", microService.ServiceName)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", NewIOException(err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var response registry.GetExistenceResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return "", NewJSONException(err, string(body))
		}
		microService.ServiceId = response.ServiceId
		return response.ServiceId, nil
	}
	if resp.StatusCode == 400 {
		return "", fmt.Errorf("client seems to have erred, error: %s", body)
	}
	return "", fmt.Errorf("RegisterService failed, MicroServiceName/responseStatusCode/responsebody: %s/%d/%s",
		microService.ServiceName, resp.StatusCode, string(body))
}

// GetProviders gets a list of provider for a particular consumer
func (c *RegistryClient) GetProviders(consumer string, opts ...CallOption) (*MicroServiceProvideResponse, error) {
	copts := &CallOptions{}
	for _, opt := range opts {
		opt(copts)
	}
	providersURL := c.formatURL(fmt.Sprintf("%s%s/%s/providers", MSAPIPath, MicroservicePath, consumer), nil, copts)
	resp, err := c.HTTPDo("GET", providersURL, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get Providers failed, error: %s, MicroServiceid: %s", err, consumer)
	}
	if resp == nil {
		return nil, fmt.Errorf("get Providers failed, response is empty, MicroServiceid: %s", consumer)
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Get Providers failed, body is empty,  error: %s, MicroServiceid: %s", err, consumer)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		p := &MicroServiceProvideResponse{}
		err = json.Unmarshal(body, p)
		if err != nil {
			return nil, err
		}
		return p, nil
	}
	return nil, fmt.Errorf("get Providers failed, MicroServiceid: %s, response StatusCode: %d, response body: %s",
		consumer, resp.StatusCode, string(body))
}

// AddSchemas adds a schema contents to the services registered in service-center
func (c *RegistryClient) AddSchemas(microServiceID, schemaName, schemaInfo string) error {
	if microServiceID == "" {
		return errors.New("invalid micro service ID")
	}

	schemaURL := c.formatURL(fmt.Sprintf("%s%s/%s%s", MSAPIPath, MicroservicePath, microServiceID, SchemaPath), nil, nil)
	h := sha256.New()
	_, err := h.Write([]byte(schemaInfo))
	if err != nil {
		return err
	}
	request := &registry.ModifySchemasRequest{
		Schemas: []*registry.Schema{{
			SchemaId: schemaName,
			Schema:   schemaInfo,
			Summary:  fmt.Sprintf("%x", h.Sum(nil))}},
	}
	body, err := json.Marshal(request)
	if err != nil {
		return NewJSONException(err, string(body))
	}

	resp, err := c.HTTPDo("POST", schemaURL, nil, body)
	if err != nil {
		return err
	}

	if resp == nil {
		return fmt.Errorf("add schemas failed, response is empty")
	}

	if resp.StatusCode != http.StatusOK {
		return NewCommonException("add micro service schema failed. response StatusCode: %d, response body: %s",
			resp.StatusCode, string(httputil.ReadBody(resp)))
	}

	return nil
}

// GetSchema gets Schema list for the microservice from service-center
func (c *RegistryClient) GetSchema(microServiceID, schemaName string, opts ...CallOption) ([]byte, error) {
	if microServiceID == "" {
		return []byte(""), errors.New("invalid micro service ID")
	}
	copts := &CallOptions{}
	for _, opt := range opts {
		opt(copts)
	}
	url := c.formatURL(fmt.Sprintf("%s%s/%s/%s/%s", MSAPIPath, MicroservicePath, microServiceID, "schemas", schemaName), nil, copts)
	resp, err := c.HTTPDo("GET", url, nil, nil)
	if err != nil {
		return []byte(""), err
	}
	if resp == nil {
		return []byte(""), fmt.Errorf("GetSchema failed, response is empty")
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), NewIOException(err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return body, nil
	}

	return []byte(""), err
}

// GetMicroServiceID gets the microserviceid by appID, serviceName and version
func (c *RegistryClient) GetMicroServiceID(appID, microServiceName, version, env string, opts ...CallOption) (string, error) {
	copts := &CallOptions{}
	for _, opt := range opts {
		opt(copts)
	}
	url := c.formatURL(MSAPIPath+ExistencePath, []URLParameter{
		{"type": "microservice"},
		{"appId": appID},
		{"serviceName": microServiceName},
		{"version": version},
		{"env": env},
	}, copts)
	resp, err := c.HTTPDo("GET", url, nil, nil)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", fmt.Errorf("GetMicroServiceID failed, response is empty, MicroServiceName: %s", microServiceName)
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", NewIOException(err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 500 {
		var response registry.GetExistenceResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return "", NewJSONException(err, string(body))
		}
		return response.ServiceId, nil
	}
	return "", fmt.Errorf("GetMicroServiceID failed, MicroService: %s@%s#%s, response StatusCode: %d, response body: %s, URL: %s",
		microServiceName, appID, version, resp.StatusCode, string(body), url)
}

// GetAllMicroServices gets list of all the microservices registered with Service-Center
func (c *RegistryClient) GetAllMicroServices(opts ...CallOption) ([]*registry.MicroService, error) {
	copts := &CallOptions{}
	for _, opt := range opts {
		opt(copts)
	}
	url := c.formatURL(MSAPIPath+MicroservicePath, nil, copts)
	resp, err := c.HTTPDo("GET", url, nil, nil)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("GetAllMicroServices failed, response is empty")
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, NewIOException(err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var response registry.GetServicesResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, NewJSONException(err, string(body))
		}
		return response.Services, nil
	}
	return nil, fmt.Errorf("GetAllMicroServices failed, response StatusCode: %d, response body: %s", resp.StatusCode, string(body))
}

// GetAllApplications returns the list of all the applications which is registered in governance-center
func (c *RegistryClient) GetAllApplications(opts ...CallOption) ([]string, error) {
	copts := &CallOptions{}
	for _, opt := range opts {
		opt(copts)
	}
	governanceURL := c.formatURL(GovernAPIPATH+AppsPath, nil, copts)
	resp, err := c.HTTPDo("GET", governanceURL, nil, nil)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("GetAllApplications failed, response is empty")
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, NewIOException(err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var response registry.GetAppsResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, NewJSONException(err, string(body))
		}
		return response.AppIds, nil
	}
	return nil, fmt.Errorf("GetAllApplications failed, response StatusCode: %d, response body: %s", resp.StatusCode, string(body))
}

// GetMicroService returns the microservices by ID
func (c *RegistryClient) GetMicroService(microServiceID string, opts ...CallOption) (*registry.MicroService, error) {
	copts := &CallOptions{}
	for _, opt := range opts {
		opt(copts)
	}
	microserviceURL := c.formatURL(fmt.Sprintf("%s%s/%s", MSAPIPath, MicroservicePath, microServiceID), nil, copts)
	resp, err := c.HTTPDo("GET", microserviceURL, nil, nil)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("GetMicroService failed, response is empty, MicroServiceId: %s", microServiceID)
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, NewIOException(err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var response registry.GetServiceResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, NewJSONException(err, string(body))
		}
		return response.Service, nil
	}
	return nil, fmt.Errorf("GetMicroService failed, MicroServiceId: %s, response StatusCode: %d, response body: %s\n, microserviceURL: %s", microServiceID, resp.StatusCode, string(body), microserviceURL)
}

//BatchFindInstances fetch instances based on service name, env, app and version
//finally it return instances grouped by service name
func (c *RegistryClient) BatchFindInstances(consumerID string, keys []*registry.FindService, opts ...CallOption) (*registry.BatchFindInstancesResponse, error) {
	copts := &CallOptions{Revision: c.revision}
	for _, opt := range opts {
		opt(copts)
	}
	if len(keys) == 0 {
		return nil, ErrEmptyCriteria
	}
	url := c.formatURL(MSAPIPath+BatchInstancePath, []URLParameter{
		{"type": "query"},
	}, copts)
	r := &registry.BatchFindInstancesRequest{
		ConsumerServiceId: consumerID,
		Services:          keys,
	}
	rBody, err := json.Marshal(r)
	if err != nil {
		return nil, NewJSONException(err, string(rBody))
	}
	resp, err := c.HTTPDo("POST", url, http.Header{"X-ConsumerId": []string{consumerID}}, rBody)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("BatchFindInstances failed, response is empty")
	}
	body := httputil.ReadBody(resp)
	if resp.StatusCode == http.StatusOK {
		var response *registry.BatchFindInstancesResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, NewJSONException(err, string(body))
		}

		return response, nil
	}
	return nil, fmt.Errorf("batch find failed, status %d, body %s", resp.StatusCode, body)
}

// FindMicroServiceInstances find microservice instance using consumerID, appID, name and version rule
func (c *RegistryClient) FindMicroServiceInstances(consumerID, appID, microServiceName,
	versionRule string, opts ...CallOption) ([]*registry.MicroServiceInstance, error) {
	copts := &CallOptions{Revision: c.revision}
	for _, opt := range opts {
		opt(copts)
	}
	microserviceInstanceURL := c.formatURL(MSAPIPath+InstancePath, []URLParameter{
		{"appId": appID},
		{"serviceName": microServiceName},
		{"version": versionRule},
	}, copts)

	resp, err := c.HTTPDo("GET", microserviceInstanceURL, http.Header{"X-ConsumerId": []string{consumerID}}, nil)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("FindMicroServiceInstances failed, response is empty, appID/MicroServiceName/version: %s/%s/%s", appID, microServiceName, versionRule)
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, NewIOException(err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var response registry.GetInstancesResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, NewJSONException(err, string(body))
		}
		r := resp.Header.Get(HeaderRevision)
		if r != c.revision && r != "" {
			c.revision = r
			openlog.Debug("service center has new revision " + c.revision)
		}

		return response.Instances, nil
	}
	if resp.StatusCode == http.StatusNotModified {
		return nil, ErrNotModified
	}
	if resp.StatusCode == http.StatusBadRequest {
		if strings.Contains(string(body), "\"errorCode\":\"400012\"") {
			return nil, ErrMicroServiceNotExists
		}
	}
	return nil, fmt.Errorf("FindMicroServiceInstances failed, appID/MicroServiceName/version: %s/%s/%s, response StatusCode: %d, response body: %s",
		appID, microServiceName, versionRule, resp.StatusCode, string(body))
}

// RegisterMicroServiceInstance registers the microservice instance to Servive-Center
func (c *RegistryClient) RegisterMicroServiceInstance(microServiceInstance *registry.MicroServiceInstance) (string, error) {
	if microServiceInstance == nil {
		return "", errors.New("invalid request parameter")
	}
	request := &registry.RegisterInstanceRequest{
		Instance: microServiceInstance,
	}
	microserviceInstanceURL := c.formatURL(fmt.Sprintf("%s%s/%s%s", MSAPIPath, MicroservicePath, microServiceInstance.ServiceId, InstancePath), nil, nil)
	body, err := json.Marshal(request)
	if err != nil {
		return "", NewJSONException(err, string(body))
	}
	resp, err := c.HTTPDo("POST", microserviceInstanceURL, nil, body)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", fmt.Errorf("register instance failed, response is empty, MicroServiceId = %s", microServiceInstance.ServiceId)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", NewIOException(err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var response *registry.RegisterInstanceResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return "", NewJSONException(err, string(body))
		}
		return response.InstanceId, nil
	}
	return "", fmt.Errorf("register instance failed, MicroServiceId: %s, response StatusCode: %d, response body: %s",
		microServiceInstance.ServiceId, resp.StatusCode, string(body))
}

// GetMicroServiceInstances queries the service-center with provider and consumer ID and returns the microservice-instance
func (c *RegistryClient) GetMicroServiceInstances(consumerID, providerID string, opts ...CallOption) ([]*registry.MicroServiceInstance, error) {
	copts := &CallOptions{}
	for _, opt := range opts {
		opt(copts)
	}
	url := c.formatURL(fmt.Sprintf("%s%s/%s%s", MSAPIPath, MicroservicePath, providerID, InstancePath), nil, copts)
	resp, err := c.HTTPDo("GET", url, http.Header{
		"X-ConsumerId": []string{consumerID},
	}, nil)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("GetMicroServiceInstances failed, response is empty, ConsumerId/ProviderId = %s%s", consumerID, providerID)
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, NewIOException(err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var response registry.GetInstancesResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, NewJSONException(err, string(body))
		}
		return response.Instances, nil
	}
	return nil, fmt.Errorf("GetMicroServiceInstances failed, ConsumerId/ProviderId: %s%s, response StatusCode: %d, response body: %s",
		consumerID, providerID, resp.StatusCode, string(body))
}

// GetAllResources retruns all the list of services, instances, providers, consumers in the service-center
func (c *RegistryClient) GetAllResources(resource string, opts ...CallOption) ([]*registry.ServiceDetail, error) {
	copts := &CallOptions{}
	for _, opt := range opts {
		opt(copts)
	}
	url := c.formatURL(GovernAPIPATH+MicroservicePath, []URLParameter{
		{"options": resource},
	}, copts)
	resp, err := c.HTTPDo("GET", url, nil, nil)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("GetAllResources failed, response is empty")
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, NewIOException(err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var response registry.GetServicesInfoResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, NewJSONException(err, string(body))
		}
		return response.AllServicesDetail, nil
	}
	return nil, fmt.Errorf("GetAllResources failed, response StatusCode: %d, response body: %s", resp.StatusCode, string(body))
}

// Health returns the list of all the endpoints of SC with their status
func (c *RegistryClient) Health() ([]*registry.MicroServiceInstance, error) {
	url := ""
	if c.apiVersion == "v4" {
		url = c.formatURL(MSAPIPath+"/health", nil, nil)
	} else {
		url = c.formatURL("/health", nil, nil)
	}

	resp, err := c.HTTPDo("GET", url, nil, nil)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("query cluster info failed, response is empty")
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, NewIOException(err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var response registry.GetInstancesResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, NewJSONException(err, string(body))
		}
		return response.Instances, nil
	}
	return nil, fmt.Errorf("query cluster info failed,  response StatusCode: %d, response body: %s",
		resp.StatusCode, string(body))
}

// Heartbeat sends the heartbeat to service-senter for particular service-instance
func (c *RegistryClient) Heartbeat(microServiceID, microServiceInstanceID string) (bool, error) {
	url := c.formatURL(fmt.Sprintf("%s%s/%s%s/%s%s", MSAPIPath, MicroservicePath, microServiceID,
		InstancePath, microServiceInstanceID, HeartbeatPath), nil, nil)
	resp, err := c.HTTPDo("PUT", url, nil, nil)
	if err != nil {
		return false, err
	}
	if resp == nil {
		return false, fmt.Errorf("heartbeat failed, response is empty, MicroServiceId/MicroServiceInstanceId: %s%s", microServiceID, microServiceInstanceID)
	}
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, NewIOException(err)
		}
		return false, NewCommonException("result: %d %s", resp.StatusCode, string(body))
	}
	return true, nil
}

// UnregisterMicroServiceInstance un-registers the microservice instance from the service-center
func (c *RegistryClient) UnregisterMicroServiceInstance(microServiceID, microServiceInstanceID string) (bool, error) {
	url := c.formatURL(fmt.Sprintf("%s%s/%s%s/%s", MSAPIPath, MicroservicePath, microServiceID,
		InstancePath, microServiceInstanceID), nil, nil)
	resp, err := c.HTTPDo("DELETE", url, nil, nil)
	if err != nil {
		return false, err
	}
	if resp == nil {
		return false, fmt.Errorf("unregister instance failed, response is empty, MicroServiceId/MicroServiceInstanceId: %s/%s", microServiceID, microServiceInstanceID)
	}
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, NewIOException(err)
		}
		return false, NewCommonException("result: %d %s", resp.StatusCode, string(body))
	}
	return true, nil
}

// UnregisterMicroService un-registers the microservice from the service-center
func (c *RegistryClient) UnregisterMicroService(microServiceID string) (bool, error) {
	url := c.formatURL(fmt.Sprintf("%s%s/%s", MSAPIPath, MicroservicePath, microServiceID), []URLParameter{
		{"force": "1"},
	}, nil)
	resp, err := c.HTTPDo("DELETE", url, nil, nil)
	if err != nil {
		return false, err
	}
	if resp == nil {
		return false, fmt.Errorf("UnregisterMicroService failed, response is empty, MicroServiceId: %s", microServiceID)
	}
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, NewIOException(err)
		}
		return false, NewCommonException("result: %d %s", resp.StatusCode, string(body))
	}
	return true, nil
}

// UpdateMicroServiceInstanceStatus updates the microservicve instance status in service-center
func (c *RegistryClient) UpdateMicroServiceInstanceStatus(microServiceID, microServiceInstanceID, status string) (bool, error) {
	url := c.formatURL(fmt.Sprintf("%s%s/%s%s/%s%s", MSAPIPath, MicroservicePath, microServiceID,
		InstancePath, microServiceInstanceID, StatusPath), []URLParameter{
		{"value": status},
	}, nil)
	resp, err := c.HTTPDo("PUT", url, nil, nil)
	if err != nil {
		return false, err
	}
	if resp == nil {
		return false, fmt.Errorf("UpdateMicroServiceInstanceStatus failed, response is empty, MicroServiceId/MicroServiceInstanceId/status: %s%s%s",
			microServiceID, microServiceInstanceID, status)
	}
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, NewIOException(err)
		}
		return false, NewCommonException("result: %d %s", resp.StatusCode, string(body))
	}
	return true, nil
}

// UpdateMicroServiceInstanceProperties updates the microserviceinstance  prooperties in the service-center
func (c *RegistryClient) UpdateMicroServiceInstanceProperties(microServiceID, microServiceInstanceID string,
	microServiceInstance *registry.MicroServiceInstance) (bool, error) {
	if microServiceInstance.Properties == nil {
		return false, errors.New("invalid request parameter")
	}
	request := registry.RegisterInstanceRequest{
		Instance: microServiceInstance,
	}
	url := c.formatURL(fmt.Sprintf("%s%s/%s%s/%s%s", MSAPIPath, MicroservicePath, microServiceID, InstancePath, microServiceInstanceID, PropertiesPath), nil, nil)
	body, err := json.Marshal(request.Instance)
	if err != nil {
		return false, NewJSONException(err, string(body))
	}

	resp, err := c.HTTPDo("PUT", url, nil, body)

	if err != nil {
		return false, err
	}
	if resp == nil {
		return false, fmt.Errorf("UpdateMicroServiceInstanceProperties failed, response is empty, MicroServiceId/microServiceInstanceID: %s/%s",
			microServiceID, microServiceInstanceID)
	}
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, NewIOException(err)
		}
		return false, NewCommonException("result: %d %s", resp.StatusCode, string(body))
	}
	return true, nil
}

// UpdateMicroServiceProperties updates the microservice properties in the servive-center
func (c *RegistryClient) UpdateMicroServiceProperties(microServiceID string, microService *registry.MicroService) (bool, error) {
	if microService.Properties == nil {
		return false, errors.New("invalid request parameter")
	}
	request := &registry.CreateServiceRequest{
		Service: microService,
	}
	url := c.formatURL(fmt.Sprintf("%s%s/%s%s", MSAPIPath, MicroservicePath, microServiceID, PropertiesPath), nil, nil)
	body, err := json.Marshal(request.Service)
	if err != nil {
		return false, NewJSONException(err, string(body))
	}

	resp, err := c.HTTPDo("PUT", url, nil, body)

	if err != nil {
		return false, err
	}
	if resp == nil {
		return false, fmt.Errorf("UpdateMicroServiceProperties failed, response is empty, MicroServiceId: %s", microServiceID)
	}
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, NewIOException(err)
		}
		return false, NewCommonException("result: %d %s", resp.StatusCode, string(body))
	}
	return true, nil
}

// Close closes the connection with Service-Center
func (c *RegistryClient) Close() error {
	for k, v := range c.conns {
		err := v.Close()
		if err != nil {
			return fmt.Errorf("error:%s, microServiceID = %s", err.Error(), k)
		}
		delete(c.conns, k)
	}
	return nil
}

// WatchMicroService creates a web socket connection to service-center to keep a watch on the providers for a micro-service
func (c *RegistryClient) WatchMicroService(microServiceID string, callback func(*MicroServiceInstanceChangedEvent)) error {
	if ready, ok := c.watchers[microServiceID]; !ok || !ready {
		c.mutex.Lock()
		if ready, ok := c.watchers[microServiceID]; !ok || !ready {
			c.watchers[microServiceID] = true
			scheme := "wss"
			if !c.Config.SSL {
				scheme = "ws"
			}
			u := url.URL{
				Scheme: scheme,
				Host:   c.getAddress(),
				Path: fmt.Sprintf("%s%s/%s%s", MSAPIPath,
					MicroservicePath, microServiceID, WatchPath),
			}
			conn, _, err := c.wsDialer.Dial(u.String(), c.GetDefaultHeaders())
			if err != nil {
				c.watchers[microServiceID] = false
				c.mutex.Unlock()
				return fmt.Errorf("watching microservice dial catch an exception,microServiceID: %s, error:%s", microServiceID, err.Error())
			}

			c.conns[microServiceID] = conn
			go func() {
				for {
					messageType, message, err := conn.ReadMessage()
					if err != nil {
						break
					}
					if messageType == websocket.TextMessage {
						var response MicroServiceInstanceChangedEvent
						err := json.Unmarshal(message, &response)
						if err != nil {
							break
						}
						callback(&response)
					}
				}
				err = conn.Close()
				if err != nil {
					openlog.Error(err.Error())
				}
				delete(c.conns, microServiceID)
				c.startBackOff(microServiceID, callback)
			}()
		}
		c.mutex.Unlock()
	}
	return nil
}

func (c *RegistryClient) getAddress() string {
	return c.pool.GetAvailableAddress()
}

func (c *RegistryClient) startBackOff(microServiceID string, callback func(*MicroServiceInstanceChangedEvent)) {
	boff := retry.GetBackOff(retry.KindExponential, 1000, 30000)
	operation := func() error {
		c.mutex.Lock()
		c.watchers[microServiceID] = false
		c.getAddress()
		c.mutex.Unlock()
		err := c.WatchMicroService(microServiceID, callback)
		if err != nil {
			return err
		}
		return nil
	}

	err := backoff.Retry(operation, boff)
	if err == nil {
		return
	}
}
