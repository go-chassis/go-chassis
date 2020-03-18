/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package restfultest

import (
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/server"
	"net/http"
	"reflect"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/go-chassis/go-chassis/core/handler"
	chassisRestful "github.com/go-chassis/go-chassis/server/restful"
)

//Container is unit test solution for rest api method
type Container struct {
	container *restful.Container
	ws        *restful.WebService
}

//New create a isolated test container,
// you can register a struct, and it will be registered to a isolated container
func New(schema interface{}, chain *handler.Chain) (*Container, error) {
	chainName := ""
	if chain != nil {
		chainName = chain.Name
		handler.ChainMap[common.Provider+chainName] = chain
	}
	c := new(Container)
	c.container = restful.NewContainer()
	c.ws = new(restful.WebService)
	routes, err := chassisRestful.GetRouteSpecs(schema)
	if err != nil {
		panic(err)
	}

	var schemaName string
	tokens := strings.Split(reflect.TypeOf(schema).String(), ".")
	if len(tokens) >= 1 {
		schemaName = tokens[len(tokens)-1]
	}
	for k := range routes {
		chassisRestful.GroupRoutePath(&routes[k], schema)
		handler, err := chassisRestful.WrapHandlerChain(&routes[k], schema, schemaName, server.Options{ChainName: chainName})
		if err != nil {
			return nil, err
		}
		if err = chassisRestful.Register2GoRestful(routes[k], c.ws, handler); err != nil {
			return nil, err
		}
	}
	c.container.Add(c.ws)
	return c, nil
}

//ServeHTTP accept native httptest, after process, response writer will write response
func (c *Container) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	c.container.ServeHTTP(resp, req)
}
