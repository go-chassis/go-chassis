package schemas

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chassis/go-chassis/server/restful"
)

// RestFulRouterA is a struct used for implementation of restfull router program
type RestFulRouterA struct {
}

// Equal is method to compare given num and slice sum
func (r *RestFulRouterA) Equal(context *restful.Context) {
	params := struct {
		Num  int   `json:"num"`
		Nums []int `json:"nums"`
	}{}
	err := context.ReadEntity(&params)
	if err != nil {
		context.Write([]byte(err.Error()))
		return
	}
	var sum int
	for _, num := range params.Nums {
		sum += num
	}
	if sum == params.Num {
		context.Write([]byte(fmt.Sprintf("version V1 : given num is equal the sum of the slice , num  : %d ,sum : %d ", params.Num, sum)))
		return
	}
	context.Write([]byte(fmt.Sprintf("version V1 : given num not  equal the sum of the slice , num  : %d ,sum : %d ", params.Num, sum)))
}

// Say is method to reply version A say some info
func (r *RestFulRouterA) Say(context *restful.Context) {
	reslut := struct {
		Name string
		Addr string
		Age  int
	}{}
	err := context.ReadEntity(&reslut)
	if err != nil {
		context.Write([]byte(err.Error()))
		return
	}
	context.Write([]byte("version V1 : " + reslut.Name + " say : he is " +
		strconv.Itoa(reslut.Age) + " years ago ,live in " + reslut.Addr))
}

// Operation is method to add two num sum
func (r *RestFulRouterA) Operation(context *restful.Context) {
	paramMap := context.ReadPathParameters()
	numString1 := paramMap["num1"]
	numString2 := paramMap["num2"]
	num1, err := strconv.Atoi(numString1)
	if err != nil {
		context.WriteHeaderAndJSON(http.StatusInternalServerError, err, "application/json")
		return
	}
	num2, err := strconv.Atoi(numString2)
	if err != nil {
		context.WriteHeaderAndJSON(http.StatusInternalServerError, err, "application/json")
		return
	}
	sum := num1 + num2
	context.Write([]byte(fmt.Sprintf("version V1 : calculate the sum of the two numbers ,the sum is %d .", sum)))
}

// Info is a method used to reply version information
func (r *RestFulRouterA) Info(context *restful.Context) {
	versionInfo := struct {
		Name     string `json:"name"`
		Version  string `json:"version"`
		HostName string `json:"hostName"`
	}{
		Name:     "CHASSIS_SERVER_V1",
		Version:  "1.0",
		HostName: context.ReadRequest().Host,
	}

	context.WriteJSON(versionInfo, "application/json")
}

// URLPatterns helps to respond for corresponding API calls
func (r *RestFulRouterA) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/info", ResourceFunc: r.Info},
		{Method: http.MethodGet, Path: "/operation/{num1}/{num2}", ResourceFunc: r.Operation},
		{Method: http.MethodPost, Path: "/say", ResourceFunc: r.Say},
		{Method: http.MethodPost, Path: "/equal", ResourceFunc: r.Equal},
	}
}
