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

	"github.com/go-chassis/go-chassis/core/common"
	wp "github.com/go-chassis/go-chassis/core/router/weightpool"
	"strings"
)

type routeRuleEventListener struct{}

// update route rule of a service
func (r *routeRuleEventListener) Event(e *core.Event) {
	if e == nil {
		openlogging.Warn("Event pointer is nil")
		return
	}
	openlogging.Info("dark launch event", openlogging.WithTags(openlogging.Tags{
		"key":   e.Key,
		"event": e.EventType,
		"rule":  e.Value,
	}))
	s := strings.Replace(e.Key, DarkLaunchPrefix, "", 1)
	switch e.EventType {
	case common.Update:
		SaveRouteRule(s, e)
	case common.Create:
		SaveRouteRule(s, e)
	case common.Delete:
		cseRouter.DeleteRouteRuleByKey(s)
		openlogging.Info("route rule is removed", openlogging.WithTags(
			openlogging.Tags{
				"key": e.Key,
			}))
	}

}

//SaveRouteRule save event rule to local cache
func SaveRouteRule(service string, e *core.Event) {
	raw, ok := e.Value.(string)
	if !ok {
		openlogging.Error("invalid route rule", openlogging.WithTags(openlogging.Tags{
			"value": raw,
		}))
	}
	routeRules, err := ConvertJSON2RouteRule(raw)
	if err != nil {
		openlogging.Error("LoadRules route rule failed", openlogging.WithTags(openlogging.Tags{
			"err": err.Error(),
		}))
	}
	if router.ValidateRule(map[string][]*model.RouteRule{service: routeRules}) {
		cseRouter.SetRouteRuleByKey(service, routeRules)
		wp.GetPool().Reset(service)
		openlogging.Info("update route rule success", openlogging.WithTags(
			openlogging.Tags{
				"key":     e.Key,
				"service": service,
				"rule":    routeRules,
			}))
	}
}
