package rsrp

// Config holds configuration details for the reverse proxy
type Config struct {
	Routes []RouteRuleConfig `json:"routes"`
}

// A RouteRuleConfig is the on-disk representation of a RouteRule
type RouteRuleConfig struct {
	Match       string            `json:"match"`
	Rewrite     RewriteRuleConfig `json:"rewrite"`
	Destination string            `json:"destination"`
}

// A RewriteRuleConfig is the on-disk representation of a RewriteRule
type RewriteRuleConfig struct {
	Input  string `json:"from"`
	Output string `json:"to"`
}
