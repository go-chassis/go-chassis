package restful_test

import (
	"github.com/go-chassis/go-chassis/server/restful"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetRouteSpecs(t *testing.T) {
	_, err := restful.GetRouteSpecs(WrongSchema{})
	assert.Error(t, err)
	_, err = restful.GetRouteSpecs(&WrongSchema{})
	assert.Error(t, err)
	_, err = restful.GetRouteSpecs(&WrongSchema2{})
	assert.Error(t, err)
}

type WrongSchema struct {
}

func (r *WrongSchema) Put(b *restful.Context) {
}

//URLPatterns helps to respond for corresponding API calls
func (r *WrongSchema) URLPatterns2() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/", ResourceFuncName: "Put",
			Returns: []*restful.Returns{{Code: 200}}},
	}
}

type WrongSchema2 struct {
}

//URLPatterns helps to respond for corresponding API calls
func (r *WrongSchema2) URLPatterns() {
}
