# API Server with JSON Responses

## /ping

Any request to `/ping` will return the following response:

```
{
  "pong": true
}
```

eg `curl -X GET http://localhost:5002/ping`

## /nocontent

Any request to `/nocontent` will return a 204 status code and a message in the `X-Header` header.

eg `curl -X GET http://localhost:5002/nocontent`

## /echo

Any request to `/echo` will return a JSON payload with the following structure:

```
{
    "body": "either a string or parsed JSON",
    "headers": {
        "Header-Key": [
            "Header Value 1",
            "Header Value 2"
        ]
    },
    "method": "METHOD",
    "path": "/echo",
    "query": {
        "QueryKey": [
            "Query Value 1",
            "Query Value 2"
        ]
    }
}
```

with information from the incoming request.