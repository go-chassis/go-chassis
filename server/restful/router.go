package restful

import (
	"fmt"
	"reflect"
)

//Route describe http route path and swagger specifications for API
type Route struct {
	Method           string // Method is one of the following: GET,PUT,POST,DELETE
	Path             string // Path contains a path pattern
	ResourceFuncName string //Resource function name
}

//GetRouteSpecs is to return a rest API specification of a go struct
func GetRouteSpecs(schema interface{}) ([]Route, error) {
	rfValue := reflect.ValueOf(schema)
	name := reflect.Indirect(rfValue).Type().Name()
	urlPatternFunc := rfValue.MethodByName("URLPatterns")
	if !urlPatternFunc.IsValid() {
		return []Route{}, fmt.Errorf("<rest.RegisterResource> no 'URLPatterns' function in servant struct `%s`", name)
	}
	vals := urlPatternFunc.Call([]reflect.Value{})
	if len(vals) <= 0 {
		return []Route{}, fmt.Errorf("<rest.RegisterResource> call URLPatterns function failed in struct `%s`", name)
	}

	if !rfValue.CanInterface() {
		return []Route{}, fmt.Errorf("<rest.RegisterResource> result of 'URLPatterns' function not interface type in servant struct `%s`", name)
	}

	if routes, ok := vals[0].Interface().([]Route); ok {
		return routes, nil
	}
	return []Route{}, fmt.Errorf("<rest.RegisterResource> result of 'URLPatterns' function not []*Route type in servant struct `%s`", name)
}
