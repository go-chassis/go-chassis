// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"fmt"
	"net/url"
	"strings"
)

// URLBuilder is the string builder to build request url
type URLBuilder struct {
	Protocol      string
	Host          string
	Path          string
	URLParameters []URLParameter
	CallOptions   *CallOptions
}

func (b *URLBuilder) encodeParams(params []URLParameter) string {
	encoded := []string{}
	for _, param := range params {
		for k, v := range param {
			if k == "" || v == "" {
				continue
			}
			encoded = append(encoded, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
		}
	}
	return strings.Join(encoded, "&")
}

// String is the method to return url string
func (b *URLBuilder) String() string {
	querys := b.URLParameters
	if b.CallOptions != nil {
		if !b.CallOptions.WithoutRevision && len(b.CallOptions.Revision) > 0 {
			querys = append(querys, URLParameter{"rev": b.CallOptions.Revision})
		}
		if b.CallOptions.WithGlobal {
			querys = append(querys, URLParameter{"global": "true"})
		}
	}
	urlString := fmt.Sprintf("%s://%s%s", b.Protocol, b.Host, b.Path)
	queryString := b.encodeParams(querys)
	if len(queryString) > 0 {
		urlString += "?" + queryString
	}
	return urlString
}
