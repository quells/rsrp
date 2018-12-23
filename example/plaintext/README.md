# Plaintext / Polyglot Example Server

A HTTP server with a single route implemented two ways.

With either language

```
curl -X GET \
  http://localhost:5001/ping
```

should yield

```
HTTP/1.0 200 OK
Content-Type: text/plain; charset=utf-8

pong
```

with some additional headers.

Anything else will yield a 404 response.

## Go version

To run: `go run main.go`

## Python 2.7.x version

To run: `python server.py`
