package restful

import (
	"fmt"
	"reflect"
)

//RouteSpec is a struct
type RouteSpec struct {
	// Method is one of the following: GET,PUT,POST,DELETE
	Method string
	// Path contains a path pattern
	Path string
	//Resource function name
	ResourceFuncName string
}

//GetRouteSpecs is to return a rest API specification of a go struct
func GetRouteSpecs(schema interface{}) ([]RouteSpec, error) {
	rfValue := reflect.ValueOf(schema)
	name := reflect.Indirect(rfValue).Type().Name()
	urlPatternFunc := rfValue.MethodByName("URLPatterns")
	if !urlPatternFunc.IsValid() {
		return []RouteSpec{}, fmt.Errorf("<rest.RegisterResource> no 'URLPatterns' function in servant struct `%s`", name)
	}
	vals := urlPatternFunc.Call([]reflect.Value{})
	if len(vals) <= 0 {
		return []RouteSpec{}, fmt.Errorf("<rest.RegisterResource> call URLPatterns function failed in struct `%s`", name)
	}

	if !rfValue.CanInterface() {
		return []RouteSpec{}, fmt.Errorf("<rest.RegisterResource> result of 'URLPatterns' function not interface type in servant struct `%s`", name)
	}

	if routes, ok := vals[0].Interface().([]RouteSpec); ok {
		return routes, nil
	}
	return []RouteSpec{}, fmt.Errorf("<rest.RegisterResource> result of 'URLPatterns' function not []*RouteSpec type in servant struct `%s`", name)
}
