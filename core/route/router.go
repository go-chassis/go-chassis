package router

import (
	"errors"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/registry"
	"regexp"
	"strconv"
	"sync"
)

var dests map[string][]*config.RouteRule

var templates map[string]*config.Match

var lock sync.RWMutex

var invokeCount = initMap()

// SafeMap safe map structure
type SafeMap struct {
	sync.RWMutex
	Map map[string]int
}

// initMap initialize map
func initMap() *SafeMap {
	sm := new(SafeMap)
	sm.Map = make(map[string]int)
	return sm
}

func (sm *SafeMap) get(key string) (int, bool) {
	sm.RLock()
	value, ok := sm.Map[key]
	sm.RUnlock()
	return value, ok
}

func (sm *SafeMap) set(key string, value int) {
	sm.Lock()
	sm.Map[key] = value
	sm.Unlock()
}

// Init sets route rule, and templates
func Init(nr map[string][]*config.RouteRule, t map[string]*config.Match) {
	SetRouteRule(nr)
	SetTemplates(t)
}

// SetRouteRule set route rule
func SetRouteRule(nr map[string][]*config.RouteRule) {
	lock.Lock()
	dests = nr
	lock.Unlock()
}

// GetRouteRule get route rule
func GetRouteRule() map[string][]*config.RouteRule {
	lock.RLock()
	defer lock.RUnlock()
	return dests
}

// SetTemplates set templates
func SetTemplates(t map[string]*config.Match) {
	lock.Lock()
	templates = t
	lock.Unlock()
}

// GetTemplates get templates
func GetTemplates() map[string]*config.Match {
	lock.RLock()
	defer lock.RUnlock()
	return templates
}

// Route route the APIs
func Route(header map[string]string, si *registry.SourceInfo, inv *invocation.Invocation) error {
	rules := SortRules(inv.MicroServiceName)
	for _, rule := range rules {
		if Match(rule.Match, header, si) {
			tag, _ := FitRate(rule.Routes, inv.MicroServiceName)
			inv.Version = tag.Tags[common.BuildinTagVersion]
			if tag.Tags[common.BuildinTagApp] != "" {
				inv.AppID = tag.Tags[common.BuildinTagApp]
			}
			break
		}
	}
	if inv.Version == "" && len(rules) != 0 {
		return errors.New("There is no matched route rule! Please make sure your request is valid")
	}
	//Finally, must set app and version for a destination,
	//because sc need those, But user don't need to care, if they don't want(means don't need to write any route rule configs)
	//in server side discovery, kubernetes pod labels must be also empty
	if inv.AppID == "" {
		if si != nil {
			inv.AppID = si.Tags[common.BuildinTagApp]
		}
		if inv.AppID == "" {
			inv.AppID = common.DefaultApp
		}
	}
	if inv.Version == "" {
		inv.Version = common.DefaultVersion
	}
	return nil
}

// FitRate fit rate
func FitRate(tags []*config.RouteTag, dest string) (tag *config.RouteTag, err error) {
	if tags[0].Weight == 100 || len(tags) == 1 {
		tag = tags[0]
		return tag, nil
	}

	totalKey := dest + "-t-" + tags[0].Tags[common.BuildinTagVersion] + "-" + tags[0].Tags[common.BuildinTagApp]
	firstKey := dest + "-" + tags[0].Tags[common.BuildinTagVersion] + "-" + tags[0].Tags[common.BuildinTagApp]
	total, ok := invokeCount.get(totalKey)
	// invoke request num for dest is 0
	if !ok {
		total = 0
		invokeCount.set(firstKey, 0)
	}

	invokeCount.set(totalKey, total+1)
	// first request or only contain one rule tag, route to tags[0]
	if total == 0 {
		tag = tags[0]
		invokeCount.set(firstKey, 1)
		return tag, nil
	}

	for _, t := range tags {
		key := dest + "-" + t.Tags[common.BuildinTagVersion] + "-" + t.Tags[common.BuildinTagApp]
		percent, exist := invokeCount.get(key)
		if !exist {
			percent = 0
		}
		//currently, t does not get enough requests, then route this one to t
		if (percent * 100 / total) <= t.Weight {
			tag = t
			invokeCount.set(key, percent+1)
			break
		}
	}
	// weight is not set, then tag can not be filtered out.
	if tag == nil && len(tags) > 0 {
		tag = tags[0]
	}
	return tag, nil
}

// Match check the route rule
func Match(match config.Match, headers map[string]string, source *registry.SourceInfo) bool {
	//validate template first
	if refer := match.Refer; refer != "" {
		return SourceMatch(templates[refer], headers, source)
	}
	//match rule is not set
	if match.Source == "" && match.HTTPHeaders == nil && match.Headers == nil {
		return true
	}

	return SourceMatch(&match, headers, source)
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
	header := headers[k]
	if regex, ok := v["regex"]; ok {
		reg := regexp.MustCompilePOSIX(regex)
		if !reg.Match([]byte(header)) {
			return false
		}
		return true
	}
	if exact, ok := v["exact"]; ok {
		if exact != header {
			return false
		}
		return true
	}
	if noEqu, ok := v["noEqu"]; ok {
		if noEqu == header {
			return false
		}
		return true
	}

	headerInt, err := strconv.Atoi(header)
	if err != nil {
		return false
	}
	if noLess, ok := v["noLess"]; ok {
		head, _ := strconv.Atoi(noLess)
		if head > headerInt {
			return false
		}
		return true
	}
	if noGreater, ok := v["noGreater"]; ok {
		head, _ := strconv.Atoi(noGreater)
		if head < headerInt {
			return false
		}
		return true
	}
	if greater, ok := v["greater"]; ok {
		head, _ := strconv.Atoi(greater)
		if head >= headerInt {
			return false
		}
		return true
	}
	if less, ok := v["less"]; ok {
		head, _ := strconv.Atoi(less)
		if head <= headerInt {
			return false
		}
	}
	return true
}

// SortRules sort route rules
func SortRules(name string) []*config.RouteRule {
	slice := dests[name]
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
