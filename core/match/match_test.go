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

package match_test

import (
	"context"
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/match"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestMatch(t *testing.T) {
	b, _ := match.Match("exact", "a", "a")
	assert.True(t, b)

	match.Install("notEq", func(v, e string) bool {
		return !(v == e)
	})

	b, _ = match.Match("notEq", "a", "a")
	assert.False(t, b)
}

func TestSaveMatchPolicy(t *testing.T) {
	testName := "match-user-json"
	testMatchPolic := `
        headers:
          cookie:
            regex: "^(.*?;)?(user=jason)(;.*)?$"
          user:
            equal: jason
        apiPath:
          contains: "some/api"
        method: GET
	`
	match.SaveMatchPolicy(testMatchPolic, "servicecomb.match."+testName, testName)
}

func TestMark(t *testing.T) {

	t.Run("test match header", func(t *testing.T) {
		testName := "match-user-json-header"
		testMatchPolic := `
        headers:
          user:
            exact: jason`
		match.SaveMatchPolicy(testMatchPolic, "servicecomb.match."+testName, testName)
		i := createInvoker(map[string]string{
			"user": "jason",
		}, http.MethodPost, "")
		match.Mark(i)
		assert.Equal(t, testName, i.GetMark())
	})
	t.Run("test match method", func(t *testing.T) {
		testName := "match-user-json-method"
		testMatchPolic := `
        method: GET`
		match.SaveMatchPolicy(testMatchPolic, "servicecomb.match."+testName, testName)
		i := createInvoker(nil, http.MethodGet, "")
		match.Mark(i)
		assert.Equal(t, testName, i.GetMark())
	})

	t.Run("test match apipath", func(t *testing.T) {
		testName := "match-user-json-apipath"
		testMatchPolic := `
        apiPath: 
          contains: "path/test"`
		match.SaveMatchPolicy(testMatchPolic, "servicecomb.match."+testName, testName)
		i := createInvoker(nil, http.MethodPost, "cse://127.0.0.1:9992/path/test")
		match.Mark(i)
		assert.Equal(t, testName, i.GetMark())
	})
}

func createInvoker(headers map[string]string, method, url string) *invocation.Invocation {
	i := invocation.New(context.Background())
	i.Metadata = make(map[string]interface{})
	for k, v := range headers {
		i.SetHeader(k, v)
	}
	i.Args, _ = rest.NewRequest(method, url, nil)
	return i
}
