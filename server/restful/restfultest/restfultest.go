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
	"net/http"
	"reflect"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	chassisRestful "github.com/go-chassis/go-chassis/server/restful"
	"github.com/go-mesh/openlogging"
)

//Container is unit test solution for rest api method
type Container struct {
	container *restful.Container
	ws        *restful.WebService
}

//New create a isolated test container,
// you can register a struct, and it will be registered to a isolated container
func New(schema interface{}, chain *handler.Chain) (*Container, error) {
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
		handleFunc, err := chassisRestful.BuildRouteHandler(&routes[k], schema)

		handler := func(req *restful.Request, rep *restful.Response) {
			defer func() {
				if r := recover(); r != nil {
					openlogging.Error(fmt.Sprintf("handle request panic. path:%s, panic:%s", routes[k].Path, r))
					if err := rep.WriteErrorString(http.StatusInternalServerError, "server got a panic, plz check log."); err != nil {
						openlogging.Error("write response failed when handler panic, err:" + err.Error())
					}
				}
			}()
			inv, err := chassisRestful.HTTPRequest2Invocation(req, schemaName, routes[k].ResourceFuncName)
			if err != nil {
				openlogging.Error("transfer http request to invocation failed", openlogging.WithTags(openlogging.Tags{
					"err": err.Error(),
				}))
				return
			}

			if chain != nil {
				chain.Next(inv, func(ir *invocation.Response) error {
					if ir.Err != nil {
						if rep != nil {
							rep.WriteHeader(ir.Status)
						}
						return ir.Err
					}
					chassisRestful.Invocation2HTTPRequest(inv, req)
					return nil
				})
			}

			bs := chassisRestful.NewBaseServer(context.Background())
			bs.Req = req
			bs.Resp = rep

			handleFunc(bs)
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
