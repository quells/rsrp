# Really Simple Reverse Proxy

A configuration-based tool to route HTTP(S) traffic to other servers/services based on the URL path of incoming requests.

Note that this is best suited for API servers; HTML becomes more difficult because the paths change.

## Example

See [the example project](tree/master/example) with a small HTTP server using `rsrp` to proxy two other simple HTTP services.

## Configuration Specification

JSON is used for the default configuration format. Other formats can be supported (see the types in `rules.go`).

```
{
  "routes": [
    {
      "match": "regex",
      "rewrite": {
        "from": "regex with capturing groups",
        "to": "path with $"
      },
      "destination": "scheme, hostname, and (optional) port"
    }
  ]
}
```

An array of routes are expected. Routes are queried in the order specified in this array. Only the first matching route is used.

Each route has a `match` field which will match against the incoming request URL path.

Each route has a rewrite rule.

A rewrite rule has a `from` field which will capture parts of the incoming request URL path.

A rewrite rule has a `to` field which will build up the proxied URL path. Capture groups are specified using 1-indexed $ syntax.

Each route has a destination, which must have the scheme and hostname of the destination server/service. A port can also be specified.

## Version History

### 1.0 (2018-12-22)

- MVP reverse proxy with simple regex-based rewrite rules
