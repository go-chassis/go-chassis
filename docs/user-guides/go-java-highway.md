### Communication between GO consumer and JAVA provider using highway protocol:

#### GO consumer
   Go consumer uses invoker.Invoker() call to make highway communication
   ![Go consumer](images/AddAndShowEmployee.png?raw=true "Go consumer")
    
```    
   Parameters of Invoke:
    1) Context
    2) MicroserviceName
    3) SchemaID
    4) operationID
    5) Input argument
    6) Response argument
```    
   In the employ.bp.go file the structure EmployeeStruct is been used as the input argument and Response argument
   
   ![EmployeeStruct](images/EmployeeStruct.png?raw=true "EmployeeStruct")
    
   #### Java provider:
   
   Microservicename is the name provider in the microservice.yaml file . In this example it is "springboot".
   
   ![microservice.yaml](images/microservice.png?raw=true "microservice.yaml")
   
   SchemaId is the schemaID defined in the java provider. In this example it is "hello".
   OperationId is the OperationName in the java provider. In this example it is  "addAndShowEmploy".
   
   ![helloservice](images/helloservice.png?raw=true "helloservice")
   
   Employ class which has the member variables "Name" and "Phone" is used as input parameter for the operation and also response    for this api.
   
   ![Employ](images/Employ.png?raw=true "Employ")
    
