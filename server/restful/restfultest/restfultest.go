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
	"github.com/emicklei/go-restful"
	chassisRestful "github.com/go-chassis/go-chassis/server/restful"
	"net/http"
	"net/http/httptest"
)

//NewRestfulContext is a function which return context for http unit test
//it will leverage httptest package to new a response recorder
func NewRestfulContext(ctx context.Context, req *http.Request) *chassisRestful.Context {
	return &chassisRestful.Context{
		Req:  restful.NewRequest(req),
		Resp: restful.NewResponse(httptest.NewRecorder()),
		Ctx:  ctx,
	}
}
