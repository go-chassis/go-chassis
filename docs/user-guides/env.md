# Supported Environments
This shows environments that go-chassis supports, 
you can control the value by setting environment

| name   |      description      |  
|----------|:-------------:|
|HOSTING_SERVER_IP |  the IP of a VM, will be added into the metadata of instance |
|SCHEMA_ROOT |    where to read the schema files, the path format is {SCHEMA_ROOT}/{SERVICE_NAME}/{schema_id}.yaml  | 
|PAAS_CSE_ENDPOINT | address of config center and service center |
|CSE_REGISTRY_ADDR  | address of service center|
|CSE_CONFIG_CENTER_ADDR  | address of config center|
|CSE_MONITOR_SERVER_ADDR  | address of dashboard service|
|SERVICE_NAME|service name will be registered in service center|
|VERSION|version will be registered in service center|
