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

package monitoring_test

import (
	"errors"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/pkg/metrics"
	"github.com/prometheus/common/expfmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/handler"
	_ "github.com/go-chassis/go-chassis/v2/middleware/monitoring"
	"github.com/go-chassis/go-chassis/v2/server/restful"
	"github.com/go-chassis/go-chassis/v2/server/restful/restfultest"
	"github.com/stretchr/testify/assert"
)

type DummyResource struct {
}

func (r *DummyResource) Sayhello(b *restful.Context) {
	id := b.ReadPathParameter("userid")
	b.Write([]byte(id))
}

func (r *DummyResource) Err(b *restful.Context) {
	b.WriteError(500, errors.New("err"))
}

// URLPatterns helps to respond for corresponding API calls
func (r *DummyResource) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/sayhello/{userid}", ResourceFunc: r.Sayhello,
			Returns: []*restful.Returns{{Code: 200}}},
		{Method: http.MethodGet, Path: "/err", ResourceFunc: r.Err,
			Returns: []*restful.Returns{{Code: 200}}},
	}
}

func TestNewWithChain(t *testing.T) {
	archaius.Init(archaius.WithMemorySource())
	archaius.Set("servicecomb.metrics.enableGoRuntimeMetrics", false)
	metrics.Init()
	r, _ := http.NewRequest("GET", "/sayhello/some_user", nil)
	chain, err := handler.CreateChain(common.Provider, "testChain", "monitoring")
	assert.NoError(t, err)
	assert.Equal(t, "", r.Header.Get("test"))
	c, err := restfultest.New(&DummyResource{}, chain)
	assert.NoError(t, err)
	resp := httptest.NewRecorder()
	c.ServeHTTP(resp, r)
	c.ServeHTTP(resp, r)

	resp2 := httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/err", nil)
	c.ServeHTTP(resp2, r)
	body, err := io.ReadAll(resp2.Body)
	assert.NoError(t, err)
	assert.Equal(t, "err", string(body))
	assert.Equal(t, http.StatusInternalServerError, resp2.Code)
	c.ServeHTTP(resp2, r)

	mfs, err := metrics.GetSystemPrometheusRegistry().Gather()
	assert.NoError(t, err)
	w := io.Writer(os.Stdout)

	enc := expfmt.NewEncoder(w, expfmt.FmtText)

	var urlList []string
	urlList = append(urlList, "/err")
	urlList = append(urlList, "/sayhello/{userid}")

	for _, mf := range mfs {
		err := enc.Encode(mf)
		assert.NoError(t, err)
		for _, metric := range mf.Metric {
			for _, label := range metric.Label {
				if *label.Name == "API" {
					assert.Equal(t, true, checkContains(urlList, *label.Value))
				}
			}
		}
	}
}

func checkContains(urlList []string, url string) bool {
	for _, rulPath := range urlList {
		if rulPath == url {
			return true
		}
	}
	return false
}
