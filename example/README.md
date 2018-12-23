# Example Proxy with rsrp

There are two services, `plaintext` and `json`. We would like to access both from the same hostname.

We will create a simple server (`main.go`) to route requests to the appropriate service based on the request's URL path.

Any path which starts with `/plaintext` will be routed to the `plaintext` service running on port :5001. Any path which starts with `/json` will be routed to the `json` service running on port :5002.

These rules are specified in `config.json`.

The prefixes will be stripped out of the re-routed requests, so the `plaintext` and `json` services don't need to know what their relative path is - as far as they know, they are at the root.
