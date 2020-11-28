package servicecomb

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/openlog"
)

//ConvertJSON2RouteRule parse raw json from cse server to route rule config
func ConvertJSON2RouteRule(raw string) ([]*config.RouteRule, error) {
	rule := &config.DarkLaunchRule{}
	if err := json.Unmarshal([]byte(raw), rule); err != nil {
		return nil, err
	}
	routeRules := DarkLaunchRule2RouteRule(rule)
	return routeRules, nil
}

// DarkLaunchRule2RouteRule translates dark launch rule to route rule
func DarkLaunchRule2RouteRule(rule *config.DarkLaunchRule) []*config.RouteRule {
	if rule.Type == DarkLaunchTypeRate {
		routes := make([]*config.RouteTag, 0)
		for _, v := range rule.Items {
			weight, _ := strconv.Atoi(v.PolicyCondition)
			version := strings.Replace(v.GroupCondition, "version=", "", 1)

			newTag := generateRouteTags(weight, strings.Split(version, ","))
			routes = append(routes, newTag...)

		}
		return []*config.RouteRule{{
			Routes:     routes,
			Precedence: 1,
		}}
	}
	if rule.Type == DarkLaunchTypeRule {
		rules := make([]*config.RouteRule, len(rule.Items))
		for i, v := range rule.Items {
			con := v.PolicyCondition
			version := strings.Replace(v.GroupCondition, "version=", "", 1)
			match := config.Match{
				HTTPHeaders: map[string]map[string]string{},
				Headers:     map[string]map[string]string{},
			}

			if strings.Contains(con, "!=") {
				setHeadersAndHTTPHeaders(&match, v.CaseInsensitive, "noEqu", con, "!=")
			} else if strings.Contains(con, ">=") {
				setHeadersAndHTTPHeaders(&match, v.CaseInsensitive, "noLess", con, ">=")
			} else if strings.Contains(con, "<=") {
				setHeadersAndHTTPHeaders(&match, v.CaseInsensitive, "noGreater", con, "<=")
			} else if strings.Contains(con, "=") {
				setHeadersAndHTTPHeaders(&match, v.CaseInsensitive, "exact", con, "=")
			} else if strings.Contains(con, ">") {
				setHeadersAndHTTPHeaders(&match, v.CaseInsensitive, "greater", con, ">")
			} else if strings.Contains(con, "<") {
				setHeadersAndHTTPHeaders(&match, v.CaseInsensitive, "less", con, "<")
			} else if strings.Contains(con, "~") {
				setHeadersAndHTTPHeaders(&match, v.CaseInsensitive, "regex", con, "~")
			}
			newRule := &config.RouteRule{
				Routes:     generateRouteTags(100, strings.Split(version, ",")),
				Match:      match,
				Precedence: 1,
			}
			rules[i] = newRule
		}
		return rules
	}
	return nil
}

// generateRouteTags generate route tags
func generateRouteTags(weights int, versions []string) []*config.RouteTag {
	length := len(versions)
	if length == 1 {
		return []*config.RouteTag{{
			Weight: weights,
			Tags:   map[string]string{"version": versions[0]},
		}}
	}

	tags := make([]*config.RouteTag, length)
	for i, v := range versions {
		tags[i] = &config.RouteTag{
			Weight: weights / length,
			Tags:   map[string]string{"version": v},
		}
	}
	return tags
}
func caseInsensitiveToString(isCaseInsensitive bool) string {
	if isCaseInsensitive {
		return common.TRUE
	}
	return common.FALSE
}
func setHeadersAndHTTPHeaders(match *config.Match, isCaseInsensitive bool, cKey, con, sp string) {
	cons := strings.Split(con, sp)
	if len(cons) != 2 {
		openlog.Error(fmt.Sprintf("set router conf to headers failed , conf : %s", con))
		return
	}
	pkey := toCamelCase(cons[0])
	(*match).HTTPHeaders[pkey] = map[string]string{
		cKey:              cons[1],
		"caseInsensitive": caseInsensitiveToString(isCaseInsensitive),
	}
	(*match).Headers[pkey] = map[string]string{
		cKey:              cons[1],
		"caseInsensitive": caseInsensitiveToString(isCaseInsensitive),
	}

}
func toCamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Title(s)
	return strings.Replace(s, " ", "", -1)
}
