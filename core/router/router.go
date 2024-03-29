// Package router expose API for user to get or set route rule
package router

import (
	"errors"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/marker"
	"strings"

	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	wp "github.com/go-chassis/go-chassis/v2/core/router/weightpool"
	"github.com/go-chassis/openlog"
)

// Router return route rule, you can also set custom route rule
type Router interface {
	Init(Options) error
	SetRouteRule(map[string][]*config.RouteRule)
	FetchRouteRuleByServiceName(service string) []*config.RouteRule
	ListRouteRule() map[string][]*config.RouteRule
}

// ErrNoExist means if there is no router implementation
var ErrNoExist = errors.New("router not exists")
var routerServices = make(map[string]func() (Router, error))

// DefaultRouter is current router implementation
var DefaultRouter Router

// InstallRouterPlugin install router plugin
func InstallRouterPlugin(name string, f func() (Router, error)) {
	openlog.Info("installed route rule plugin: " + name)
	routerServices[name] = f
}

// BuildRouter create a router
func BuildRouter(name string) error {
	f, ok := routerServices[name]
	if !ok {
		return ErrNoExist
	}
	r, err := f()
	if err != nil {
		return err
	}
	DefaultRouter = r
	return nil
}

// Route decide the target service metadata
// it decide based on configuration of route rule
// it will set RouteTag to invocation
func Route(header map[string]string, si *registry.SourceInfo, inv *invocation.Invocation) error {
	rules := SortRules(inv.MicroServiceName)
	for _, rule := range rules {
		if Match(inv, rule.Match, header, si) {
			tag := FitRate(rule.Routes, GenWeightPoolKey(inv.MicroServiceName, rule.Precedence))
			inv.RouteTags = routeTagToTags(tag)
			break
		}
	}
	return nil
}

// FitRate fit rate
func FitRate(tags []*config.RouteTag, dest string) *config.RouteTag {
	if tags[0].Weight == 100 {
		return tags[0]
	}

	pool, ok := wp.GetPool().Get(dest)
	if !ok {
		pool = wp.NewPool(tags...)
		wp.GetPool().Set(dest, pool)
	}
	return pool.PickOne()
}

// match check the route rule
func Match(inv *invocation.Invocation, matchConf config.Match, headers map[string]string, source *registry.SourceInfo) bool {
	//validate template first
	if refer := matchConf.Refer; refer != "" {
		marker.Mark(inv)
		return inv.GetMark() == matchConf.Refer
	}
	//matchConf rule is not set
	if matchConf.Source == "" && matchConf.HTTPHeaders == nil && matchConf.Headers == nil {
		return true
	}

	return SourceMatch(&matchConf, headers, source)
}

// SourceMatch check the source route
func SourceMatch(match *config.Match, headers map[string]string, source *registry.SourceInfo) bool {
	//source not match
	if match.Source != "" && match.Source != source.Name {
		return false
	}
	//source tags not match
	if len(match.SourceTags) != 0 {
		for k, v := range match.SourceTags {
			if v != source.Tags[k] {
				return false
			}
		}
	}

	//source headers not match
	if match.Headers != nil {
		for k, v := range match.Headers {
			if !isMatch(headers, k, v) {
				return false
			}
			continue
		}
	}
	if match.HTTPHeaders != nil {
		for k, v := range match.HTTPHeaders {
			if !isMatch(headers, k, v) {
				return false
			}
			continue
		}
	}
	return true
}

// isMatch check the route rule
func isMatch(headers map[string]string, k string, v map[string]string) bool {
	header := valueToUpper(v["caseInsensitive"], headers[k])
	for op, exp := range v {
		if op == "caseInsensitive" {
			continue
		}
		if ok, err := marker.Match(op, header, valueToUpper(v["caseInsensitive"], exp)); !ok || err != nil {
			return false
		}
	}
	return true
}

func valueToUpper(b, value string) string {
	if b == common.TRUE {
		value = strings.ToUpper(value)
	}

	return value
}

// SortRules sort route rules
func SortRules(name string) []*config.RouteRule {
	if DefaultRouter == nil {
		openlog.Debug("router not available")
	}
	slice := DefaultRouter.FetchRouteRuleByServiceName(name)
	return QuickSort(0, len(slice)-1, slice)
}

// QuickSort for sorting the routes it will follow quicksort technique
func QuickSort(left int, right int, rules []*config.RouteRule) (s []*config.RouteRule) {
	s = rules
	if left >= right {
		return
	}

	i := left
	j := right
	base := s[left]
	var tmp *config.RouteRule
	for i != j {
		for s[j].Precedence <= base.Precedence && i < j {
			j--
		}
		for s[i].Precedence >= base.Precedence && i < j {
			i++
		}
		if i < j {
			tmp = s[i]
			s[i] = s[j]
			s[j] = tmp
		}
	}
	//move base to the current position of i&j
	s[left] = s[i]
	s[i] = base

	QuickSort(left, i-1, s)
	QuickSort(i+1, right, s)

	return
}
