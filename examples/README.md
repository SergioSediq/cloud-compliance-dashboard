# Examples — compliance viewer

With the server listening on `:9090` (`go run ./cmd/server` from this module’s root):

```bash
curl -sS http://127.0.0.1:9090/api/v1/version | jq .
curl -sS "http://127.0.0.1:9090/api/v1/checks?tag=CIS" | jq .
curl -sS http://127.0.0.1:9090/api/v1/checks/cis-1.1 | jq .
curl -sS -o checks.csv "http://127.0.0.1:9090/api/v1/export/checks.csv"
```

`CHECKS_PATH` points at the JSON fixture; override to point at your own file when experimenting.
