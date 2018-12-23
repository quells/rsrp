# Example Proxy with rsrp

There are two services, `plaintext` and `json`. We would like to access both from the same hostname.

We will create a simple server (`main.go`) to route requests to the appropriate service based on the request's URL path.

Any path which starts with `/plaintext` will be routed to the `plaintext` service running on port :5001. Any path which starts with `/json` will be routed to the `json` service running on port :5002.

These rules are specified in `config.json`.

The prefixes will be stripped out of the re-routed requests, so the `plaintext` and `json` services don't need to know what their relative path is - as far as they know, they are at the root.

## To Run

Open 4 terminals.

In terminal 1:

```bash
$ cd json
$ go run main.go
```

In terminal 2:

```bash
$ cd plaintext
$ go run main.go
$ # or use Python 2.7.x - there is no language requirement for a proxied server so long as it speaks HTTP
$ python server.py
```

In terminal 3:

```bash
$ go run main.go config.json
```

In terminal 4:

```bash
# Plaintext ping/pong response
$ curl -X GET \
    http://localhost:5000/plaintext/ping

# JSON ping/pong response
$ curl -X GET \
    http://localhost:5000/json/ping

# Response with a non-200 status code
$ curl -X GET \
    http://localhost:5000/json/nocontent

# Response with a JSON payload describing the incoming request
$ curl -X GET \
    http://localhost:5000/json/echo
# Note the path in the response - as far as the service knows, it received a request for /echo
```

## Troubleshooting

If one of the services won't start, you might already have something running on that port. Either quit the other service or change all occurrences of that port number to a port which is not currently occupied.
