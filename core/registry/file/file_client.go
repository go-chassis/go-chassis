package file

import (
	"encoding/json"
	"fmt"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/util/fileutil"
	"github.com/ServiceComb/go-sc-client/model"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	// ServiceJSON service json
	ServiceJSON = "service.json"
)

type localFileData struct {
	ServiceName  string `json:"serviceName,omitempty"`
	InstanceData *model.MicroServiceInstance
}

// Options struct having addresses
type Options struct {
	Addrs []string
}

type fileClient struct {
	Addresses []string
}

type serviceData struct {
	Service []*service `json:"service,omitempty"`
}
type service struct {
	Name     string   `json:"name,omitempty"`
	Instance []string `json:"instance,omitempty"`
}

func (f *fileClient) Initialize(opt Options) {
	f.Addresses = opt.Addrs
}

func (f *fileClient) FindMicroServiceInstances(microServiceName string) ([]*model.MicroServiceInstance, error) {
	var instanceData []*model.MicroServiceInstance

	data := f.getInstanceDataFromFile()
	if data == nil {
		return instanceData, fmt.Errorf("failed to get instance information")
	}

	localData := &localFileData{}
	for _, value := range data.Service {
		if value.Name == microServiceName {
			insData := &model.MicroServiceInstance{
				Endpoints: value.Instance,
			}
			localData.ServiceName = value.Name
			localData.InstanceData = insData

			instanceData = append(instanceData, localData.InstanceData)
			return instanceData, nil
		}
	}

	return instanceData, nil
}

func (f *fileClient) getInstanceDataFromFile() *serviceData {
	var data *serviceData
	path := strings.Join(f.Addresses, "")

	if path == "" {
		cwd, _ := fileutil.GetWorkDir()
		path = filepath.Join(cwd, "disco", ServiceJSON)
	}

	file, err := os.Open(path)
	if err != nil {
		lager.Logger.Warnf(err, "failed to open a file")
	}
	defer file.Close()

	plan, err := ioutil.ReadFile(path)
	if err != nil {
		lager.Logger.Warnf(err, "failed to do readfile operation")
	}

	err = json.Unmarshal(plan, &data)
	if err != nil {
		lager.Logger.Warnf(err, "failed to do unmarshall")
	}

	return data
}
