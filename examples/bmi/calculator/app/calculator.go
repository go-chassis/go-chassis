package app

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ServiceComb/go-chassis/core/registry"
	rf "github.com/ServiceComb/go-chassis/server/restful"
)

type CalculateBmi struct {
}

func (c *CalculateBmi) URLPatterns() []rf.Route {
	return []rf.Route{
		{http.MethodGet, "/bmi", "Calculate"},
	}
}

func (c *CalculateBmi) Calculate(ctx *rf.Context) {

	var height, weight, bmi float64
	var err error
	result := struct {
		Result     float64 `json:"result"`
		InstanceId string  `json:"instanceId"`
		CallTime   string  `json:"callTime"`
	}{}
	errorResponse := struct {
		Error string `json:"error"`
	}{}

	heightStr := ctx.ReadQueryParameter("height")
	weightStr := ctx.ReadQueryParameter("weight")

	if height, err = strconv.ParseFloat(heightStr, 10); err != nil {
		errorResponse.Error = err.Error()
		ctx.WriteHeaderAndJSON(http.StatusBadRequest, errorResponse, "application/json")
		return
	}
	if weight, err = strconv.ParseFloat(weightStr, 10); err != nil {
		errorResponse.Error = err.Error()
		ctx.WriteHeaderAndJSON(http.StatusBadRequest, errorResponse, "application/json")
		return
	}

	if bmi, err = c.BMIIndex(height, weight); err != nil {
		errorResponse.Error = err.Error()
		ctx.WriteHeaderAndJSON(http.StatusBadRequest, errorResponse, "application/json")
		return
	}

	result.Result = bmi
	//Get InstanceID
	items := registry.SelfInstancesCache.Items()
	for microserviceID := range items {
		instanceID, exist := registry.SelfInstancesCache.Get(microserviceID)
		if exist {
			result.InstanceId = instanceID.([]string)[0]
		}
	}
	result.CallTime = time.Now().String()
	ctx.WriteJSON(result, "application/json")
}

func (c *CalculateBmi) BMIIndex(height, weight float64) (float64, error) {
	if height <= 0 || weight <= 0 {
		return 0, fmt.Errorf("Arugments must be above 0")
	}
	heightInMeter := height / 100
	bmi := weight / (heightInMeter * heightInMeter)
	return bmi, nil
}
