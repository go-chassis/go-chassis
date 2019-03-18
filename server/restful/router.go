package restful

import (
	"fmt"
	"reflect"
)

//const for doc
const (
	Path  = "path"
	Query = "query"
)

//Route describe http route path and swagger specifications for API
type Route struct {
	Method           string        //Method is one of the following: GET,PUT,POST,DELETE. required
	Path             string        //Path contains a path pattern. required
	ResourceFuncName string        //the func this API calls. required
	FuncDesc         string        //tells what this route is all about. Optional.
	Parameters       []*Parameters //Parameters is a slice of request parameters for a single endpoint Optional.
	Returns          []*Returns    //what kind of response this API returns. Optional.
	Read             interface{}   //Read tells what resource type will be read from the request payload. Optional.
	Consumes         []string      //Consumes specifies that this WebService can consume one or more MIME types.
	Produces         []string      //Produces specifies that this WebService can produce one or more MIME types.
}

//Returns describe response doc
type Returns struct {
	Code    int // http response code
	Message string
	Model   interface{} // response body structure
}

//Parameters describe parameters in url path or query params
type Parameters struct {
	Name      string //parameter name
	DataType  string // string, int etc
	ParamType int    //restful.QueryParameterKind or restful.PathParameterKind
	Desc      string
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
