package cse

import (
	"strconv"
	"strings"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-mesh/openlogging"
)

// DarkLaunchRule2RouteRule translates dark launch rule to route rule
func DarkLaunchRule2RouteRule(rule *model.DarkLaunchRule) []*model.RouteRule {

	if rule.Type == DarkLaunchTypeRate {
		routes := make([]*model.RouteTag, 0)
		for _, v := range rule.Items {
			weight, _ := strconv.Atoi(v.PolicyCondition)
			version := strings.Replace(v.GroupCondition, "version=", "", 1)

			newTag := generateRouteTags(weight, strings.Split(version, ","))
			routes = append(routes, newTag...)

		}
		return []*model.RouteRule{{
			Routes:     routes,
			Precedence: 1,
		}}
	}
	if rule.Type == DarkLaunchTypeRule {
		rules := make([]*model.RouteRule, len(rule.Items))
		for i, v := range rule.Items {
			con := v.PolicyCondition
			version := strings.Replace(v.GroupCondition, "version=", "", 1)
			match := model.Match{
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
			newRule := &model.RouteRule{
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
func generateRouteTags(weights int, versions []string) []*model.RouteTag {
	length := len(versions)
	if length == 1 {
		return []*model.RouteTag{{
			Weight: weights,
			Tags:   map[string]string{"version": versions[0]},
		}}
	}

	tags := make([]*model.RouteTag, length)
	for i, v := range versions {
		tags[i] = &model.RouteTag{
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
func setHeadersAndHTTPHeaders(match *model.Match, isCaseInsensitive bool, cKey, con, sp string) {
	cons := strings.Split(con, sp)
	if len(cons) != 2 {
		openlogging.GetLogger().Errorf("set router conf to headers failed , conf : %s", con)
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
