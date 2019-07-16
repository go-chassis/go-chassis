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
	"context"
	"fmt"
	"github.com/emicklei/go-restful"
	chassisRestful "github.com/go-chassis/go-chassis/server/restful"
	"github.com/go-mesh/openlogging"
	"net/http"
	"reflect"
)

//Container is unit test solution for rest api method
type Container struct {
	container *restful.Container
	ws        *restful.WebService
}

//New create a isolated test container,
// you can register a struct, and it will be registered to a isolated container
func New(schema interface{}) (*Container, error) {
	c := new(Container)
	c.container = restful.NewContainer()
	c.ws = new(restful.WebService)
	routes, err := chassisRestful.GetRouteSpecs(schema)
	if err != nil {
		panic(err)
	}
	schemaType := reflect.TypeOf(schema)
	schemaValue := reflect.ValueOf(schema)
	for _, route := range routes {
		method, exist := schemaType.MethodByName(route.ResourceFuncName)
		if !exist {
			openlogging.GetLogger().Errorf("router func can not find: %s", route.ResourceFuncName)
			return nil, fmt.Errorf("router func can not find: %s", route.ResourceFuncName)
		}
		handler := func(req *restful.Request, rep *restful.Response) {

			bs := chassisRestful.NewBaseServer(context.Background())
			bs.Req = req
			bs.Resp = rep
			method.Func.Call([]reflect.Value{schemaValue, reflect.ValueOf(bs)})
		}
		if err := chassisRestful.Register2GoRestful(route, c.ws, handler); err != nil {
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
