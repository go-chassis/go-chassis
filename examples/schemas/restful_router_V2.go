package schemas

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chassis/go-chassis/v2/server/restful"
)

// RestFulRouterB is a struct used for implementation of restfull router program
type RestFulRouterB struct {
}

// Equal is method to compare given num and slice product
func (r *RestFulRouterB) Equal(context *restful.Context) {
	params := struct {
		Num  int   `json:"num"`
		Nums []int `json:"nums"`
	}{}
	err := context.ReadEntity(&params)
	if err != nil {
		context.Write([]byte(err.Error()))
		return
	}
	var product int = 1
	for _, num := range params.Nums {
		product *= num
	}
	if product == params.Num {
		context.Write([]byte(fmt.Sprintf("version  V2 : given num is equal the product of the slice , num  : %d ,product : %d ", params.Num, product)))
		return
	}
	context.Write([]byte(fmt.Sprintf("version V2 : given num not  equal the product of the slice , num  : %d ,product : %d ", params.Num, product)))
}

// Say is method to reply version B say some info
func (r *RestFulRouterB) Say(context *restful.Context) {
	reslut := struct {
		Name  string
		Addr  string
		Age   int
		Phone string
	}{}
	err := context.ReadEntity(&reslut)
	if err != nil {
		context.Write([]byte(err.Error()))
		return
	}
	if reslut.Phone == "" {
		reslut.Phone = "13800138000"
	}
	context.Write([]byte("version V2 : " + reslut.Name + " say : he is " + strconv.Itoa(reslut.Age) +
		" years ago , live in " + reslut.Addr + " ,phone is " + reslut.Phone))
}

// Operation is method to calculate  two num product
func (r *RestFulRouterB) Operation(context *restful.Context) {
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
	product := num1 * num2
	context.Write([]byte(fmt.Sprintf("version V2 : calculate the product of the two numbers ,the product is %d .", product)))
}

// Info is a method used to reply version information
func (r *RestFulRouterB) Info(context *restful.Context) {
	versionInfo := struct {
		Name     string `json:"name"`
		Version  string `json:"version"`
		HostName string `json:"hostName"`
		Now      int64  `json:"now"`
	}{
		Name:     "CHASSIS_SERVER_V2",
		Version:  "2.0",
		HostName: context.ReadRequest().Host,
		Now:      time.Now().UnixNano() / 1e6,
	}
	context.WriteJSON(versionInfo, "application/json")
}

// URLPatterns helps to respond for corresponding API calls
func (r *RestFulRouterB) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/info", ResourceFunc: r.Info},
		{Method: http.MethodGet, Path: "/operation/{num1}/{num2}", ResourceFunc: r.Operation},
		{Method: http.MethodPost, Path: "/say", ResourceFunc: r.Say},
		{Method: http.MethodPost, Path: "/equal", ResourceFunc: r.Equal},
	}
}
