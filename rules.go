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

// RewritePath converts a request path to the redirected path
func (rule RouteRule) RewritePath(path string) string {
	return rule.Rewrite.Input.ReplaceAllString(path, rule.Rewrite.Output)
}

// RewriteLocation converts a request path to the redirected location
func (rule RouteRule) RewriteLocation(path string) string {
	return rule.Destination + rule.RewritePath(path)
}

// A RewriteRule describes how to modify the path for a request
type RewriteRule struct {
	Input  *regexp.Regexp
	Output string
}

// NewRewriteRule converts a RewriteRuleConfig to a RewriteRule
func NewRewriteRule(config RewriteRuleConfig) (rule *RewriteRule, err error) {
	var input *regexp.Regexp
	input, err = regexp.Compile(config.Input)
	if err != nil {
		return
	}

	rule = &RewriteRule{
		Input:  input,
		Output: config.Output,
	}

	return
}
