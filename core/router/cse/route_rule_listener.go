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

package cse

import (
	"github.com/go-chassis/go-archaius/core"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/router"
	"github.com/go-mesh/openlogging"

	wp "github.com/go-chassis/go-chassis/core/router/weightpool"
)

type routeRuleEventListener struct{}

// update route rule of a service
func (r *routeRuleEventListener) Event(e *core.Event) {
	if e == nil {
		openlogging.Warn("Event pointer is nil")
		return
	}

	v := routeRuleMgr.GetConfigurationsByKey(e.Key)
	if v == nil {
		DeleteRouteRuleByKey(e.Key)
		openlogging.Info("route rule is removed", openlogging.WithTags(
			openlogging.Tags{
				"key": e.Key,
			}))
		return
	}
	routeRules, ok := v.([]*model.RouteRule)
	if !ok {
		openlogging.Error("value is not type []*RouteRule")
		return
	}

	if router.ValidateRule(map[string][]*model.RouteRule{e.Key: routeRules}) {
		SetRouteRuleByKey(e.Key, routeRules)
		wp.GetPool().Reset(e.Key)
		openlogging.Info("update route rule success", openlogging.WithTags(
			openlogging.Tags{
				"key": e.Key,
			}))
	}
}
