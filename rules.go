package rsrp

import (
	"regexp"
)

// ConvertRules converts RouteRuleConfigs to RouteRules
func ConvertRules(routes []RouteRuleConfig) (routeRules *[]RouteRule, err error) {
	rules := make([]RouteRule, len(routes))

	var rule *RouteRule
	for i, route := range routes {
		rule, err = NewRouteRule(route)
		if err != nil {
			return
		}

		rules[i] = *rule
	}

	routeRules = &rules

	return
}

// A RouteRule describes which paths to match, how to rewrite the request,
// and where to reroute the request
type RouteRule struct {
	Match       *regexp.Regexp
	Rewrite     RewriteRule
	Destination string
}

// NewRouteRule converts a RouteRuleConfig to a RouteRule
func NewRouteRule(config RouteRuleConfig) (rule *RouteRule, err error) {
	var match *regexp.Regexp
	match, err = regexp.Compile(config.Match)
	if err != nil {
		return
	}

	var rewrite *RewriteRule
	rewrite, err = NewRewriteRule(config.Rewrite)
	if err != nil {
		return
	}

	rule = &RouteRule{
		Match:       match,
		Rewrite:     *rewrite,
		Destination: config.Destination,
	}

	return
}

// A RewriteRule describes how to modify the path for a request
type RewriteRule struct {
	Source      *regexp.Regexp
	Destination string
}

// NewRewriteRule converts a RewriteRuleConfig to a RewriteRule
func NewRewriteRule(config RewriteRuleConfig) (rule *RewriteRule, err error) {
	var src *regexp.Regexp
	src, err = regexp.Compile(config.Source)
	if err != nil {
		return
	}

	rule = &RewriteRule{
		Source:      src,
		Destination: config.Destination,
	}

	return
}
