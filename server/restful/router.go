package restful

import (
	"fmt"
	"reflect"
)

//Route is a struct
type Route struct {
	// Method is one of the following: GET,PUT,POST,DELETE
	Method string
	// Path contains a path pattern
	Path string
	//Resource function name
	ResourceFuncName string
}

//GetRoutes is a function used to respond to corresponding API calls
func GetRoutes(schema interface{}) ([]Route, error) {
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
