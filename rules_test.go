package rsrp_test

import (
	"regexp"
	"testing"

	"github.com/quells/rsrp"
)

func TestRouteRule_RewritePath(t *testing.T) {
	rule := rsrp.RouteRule{
		Match: regexp.MustCompile("^/test.*$"),
		Rewrite: rsrp.RewriteRule{
			Input:  regexp.MustCompile("^/test(/.*)$"),
			Output: "/new$1",
		},
		Destination: "http://other",
	}

	rewritten := rule.RewritePath("/test/me")

	expected := "/new/me"
	if rewritten != expected {
		t.Fatalf("expected %s, found %s", expected, rewritten)
	}
}

func TestRouteRule_RewriteLocation(t *testing.T) {
	rule := rsrp.RouteRule{
		Match: regexp.MustCompile("^/test.*$"),
		Rewrite: rsrp.RewriteRule{
			Input:  regexp.MustCompile("^/test(/.*)$"),
			Output: "/new$1",
		},
		Destination: "http://other",
	}

	rewritten := rule.RewriteLocation("/test/me")

	expected := "http://other/new/me"
	if rewritten != expected {
		t.Fatalf("expected %s, found %s", expected, rewritten)
	}
}
