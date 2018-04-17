package pilot

import (
	"encoding/json"
	"fmt"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/http-client"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	// BaseRoot is the root path of pilot API
	BaseRoot = "/v1/registration"
	// DefaultAddr is the default endpoint of pilot-discovery
	DefaultAddr = "istio-pilot.istio-system:8080"
)

// EnvoyDSClient is the client implements istio/pilot v1 API
// See https://www.envoyproxy.io/docs/envoy/v1.6.0/api-v1/cluster_manager/sds#rest-api
type EnvoyDSClient struct {
	Options  Options
	protocol string
	client   *httpclient.URLClient
}

// Initialize is the func initialize the EnvoyDSClient
func (c *EnvoyDSClient) Initialize(options Options) (err error) {
	// copy options
	c.Options = options

	// set http protocol
	sslEnabled := options.TLSConfig != nil
	c.protocol = common.HTTPS
	if !sslEnabled {
		c.protocol = common.HTTP
	}

	// new rest client
	c.client, err = httpclient.GetURLClient(&httpclient.URLClientOption{
		SSLEnabled: sslEnabled,
		TLSConfig:  options.TLSConfig,
	})
	if err != nil {
		return err
	}
	return
}

// GetAllServices returns a list of Service registered by istio
func (c *EnvoyDSClient) GetAllServices() ([]*Service, error) {
	apiURL := c.getAddress() + BaseRoot
	resp, err := c.client.HttpDo("GET", apiURL, nil, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetAllServices failed, %s", err.Error())
	}
	if resp.StatusCode == http.StatusOK {
		var response []*Service
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, fmt.Errorf("GetAllServices failed, %s, response body: %s", err.Error(), string(body))
		}
		return response, nil
	}
	return nil, fmt.Errorf("GetAllServices failed, response StatusCode: %d, response body: %s",
		resp.StatusCode, string(body))
}

// GetServiceHosts returns Hosts using serviceName
func (c *EnvoyDSClient) GetServiceHosts(serviceName string) (*Hosts, error) {
	services, err := c.GetAllServices()
	if err != nil {
		return nil, err
	}

	var response Hosts
	for _, service := range services {
		if strings.Index(service.ServiceKey+".", serviceName+".") != 0 {
			continue
		}
		response.Hosts = append(response.Hosts, service.Hosts...)
	}
	return &response, nil
}

// Close is the function clean up client resources
func (c *EnvoyDSClient) Close() error {
	return nil
}

// getAddress contains a round robin lb to return registry address.
func (c *EnvoyDSClient) getAddress() string {
	next := RoundRobin(c.Options.Addrs)
	addr, err := next()
	if err != nil {
		return c.protocol + "://" + DefaultAddr
	}
	return c.protocol + "://" + addr
}
