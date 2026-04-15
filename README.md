# Compliance posture viewer

Author: **Sergio Sediq**

Go service that loads control rows from `data/checks.json`, exposes JSON under `/api/v1/*`, and ships a static UI from `cmd/server/static`. Wiring lives in `internal/server` (route table), `internal/config`, and `pkg/version` for link-time stamping. “Drift” here means **observed ≠ expected** on the same row—useful for interviews about configuration baselines, not a substitute for an audit.

## Run

```bash
go run ./cmd/server
# http://localhost:9090
```

`CHECKS_PATH` overrides the JSON path. `PORT` overrides `:9090`.

## Tests

```bash
go test ./...
```

After changing **`go.mod`** / **`go.sum`**, run **`go mod tidy`** with **Go 1.23+** on your `PATH`.

## API

| Method | Path | Notes |
|--------|------|--------|
| GET | `/health` | `{"status":"ok"}` |
| GET | `/api/v1/version` | build metadata (`pkg/version`) |
| GET | `/api/v1/checks` | optional `?tag=CIS` (framework tag) |
| GET | `/api/v1/checks/{id}` | single row |
| GET | `/api/v1/drift` | rows where observed ≠ expected |
| GET | `/api/v1/summary` | counts + per-tag histogram |
| GET | `/api/v1/frameworks` | sorted `{tag,count}` rows |
| GET | `/api/v1/export/checks.csv` | optional `?tag=` |

Request logging middleware wraps the mux; JSON responses are indented for quick eyeballing in a terminal.

## Docker

```bash
docker build -t compliance .
docker run --rm -p 9090:9090 compliance
```

## License

MIT — see [LICENSE](./LICENSE).
