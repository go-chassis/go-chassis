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

package restfultest_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/server/restful"
	"github.com/go-chassis/go-chassis/server/restful/restfultest"
	"github.com/stretchr/testify/assert"
)

type DummyResource struct {
}

func (r *DummyResource) GroupPath() string{
	return "/demo"
}

func (r *DummyResource) Sayhello(b *restful.Context) {
	id := b.ReadPathParameter("userid")
	b.Write([]byte(id))
}

func (r *DummyResource) Panic(b *restful.Context) {
	panic("panic msg")
}

//URLPatterns helps to respond for corresponding API calls
func (r *DummyResource) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/sayhello/{userid}", ResourceFuncName: "Sayhello",
			Returns: []*restful.Returns{{Code: 200}}},
		{Method: http.MethodGet, Path: "/sayhello2/{userid}", ResourceFunc:r.Sayhello,
			Returns: []*restful.Returns{{Code: 200}}},
		{Method: http.MethodGet, Path: "/panic", ResourceFunc:r.Panic,
			Returns: []*restful.Returns{{Code: 200}}},
	}
}

type FakeHandler struct {
}

func (fh *FakeHandler) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	i.SetHeader("test", "chain")
	r := &invocation.Response{}
	cb(r)
}

func (fh *FakeHandler) Name() string {
	return "test"
}
func newFakeHandler() handler.Handler {
	return &FakeHandler{}
}

func TestNew(t *testing.T) {
	r, _ := http.NewRequest("GET", "/demo/sayhello/some_user", nil)
	c, err := restfultest.New(&DummyResource{}, nil)
	assert.NoError(t, err)
	resp := httptest.NewRecorder()
	c.ServeHTTP(resp, r)
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "some_user", string(body))

	r, _ = http.NewRequest("GET", "/demo/sayhello2/another_user", nil)
	c.ServeHTTP(resp, r)
	body, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "another_user", string(body))

	r, _ = http.NewRequest("GET", "/demo/panic", nil)
	c.ServeHTTP(resp, r)
	body, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "server got a panic, plz check log.", string(body))
}

func TestNewWithChain(t *testing.T) {
	r, _ := http.NewRequest("GET", "/demo/sayhello/some_user", nil)
	handler.RegisterHandler("test", newFakeHandler)
	chain, _ := handler.CreateChain(common.Provider, "testChain", "test")
	assert.Equal(t, "", r.Header.Get("test"))
	c, err := restfultest.New(&DummyResource{}, chain)
	assert.NoError(t, err)
	resp := httptest.NewRecorder()
	c.ServeHTTP(resp, r)
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "some_user", string(body))
	assert.Equal(t, "chain", r.Header.Get("test"))
}
