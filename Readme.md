# KV Store

Supports following operations
* `/get/{key}` returns `{'value': '{key-value}'}` or 404 if key doesn't exist in the store
* `/set` accepts `application/json` type of data to set `key-values` for e.g.
```
{
    "message0": "this is my first message",
    "message1": "this is my second message"
}
```
* `/search` accepts query `?prefix={some-keys-prefix}` or `?suffix={some-keys-suffix}` and returns `{'keys':[{array-of-matching-keys}]}`
* `/prometheus` metrics are exported which supports `http_requests_total` (count), `response_status` (count) & `http_response_time_seconds`

Reachable from at `http://kv.vaibhavk.in`

To build the whole infrastructure, use `terraform/` directory

Kubernetes deployment configs can be found at `configs` directory
