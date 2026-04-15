# Compliance posture viewer

**Cloud compliance dashboard**

This is a Go service that loads control rows from `data/checks.json`, serves JSON under `/api/v1/*`, and ships a static UI from `cmd/server/static`. Routing and wiring live in `internal/server` (route table), `internal/config`, and `pkg/version` for link-time stamping. Drift here means observed state does not match expected state on the same row. That is useful for talking about configuration baselines in interviews; it is not a substitute for a formal audit.

**Author:** [Sergio Sediq](https://github.com/SergioSediq) · [LinkedIn](https://www.linkedin.com/in/sedyagho) · [sediqsergio@gmail.com](mailto:sediqsergio@gmail.com)

**Docs:** [Architecture](./docs/ARCHITECTURE.md) · [`go.mod`](./go.mod)

---

## Features

| Area | Description |
|------|-------------|
| Controls | Browse rows from JSON; optional `?tag=` filter (e.g. framework tag) |
| Drift | List rows where observed ≠ expected |
| Summary | Counts and per-tag histograms |
| Frameworks | Sorted tag and count pairs |
| Export | `GET /api/v1/export/checks.csv` with optional `?tag=` |
| Version | Build metadata via `pkg/version` |
| UI | Static HTML/CSS/JS embedded with the binary |
| Ops | Request logging middleware; JSON indented for easy `curl` review |

## Stack

| Layer | Technology |
|-------|------------|
| Language | Go 1.23+ ([`go.mod`](./go.mod)) |
| HTTP | `net/http`, Go 1.22+ `ServeMux`-style routing |
| UI | Static assets under `cmd/server/static` |
| Config | `CHECKS_PATH`, `PORT` |
| Tests | `httptest`, handler tests (including CSV and drift) |

## Repository layout

```
cloud-compliance-dashboard/
├── cmd/server/
│   ├── main.go              # timeouts, graceful shutdown
│   └── static/              # embedded UI (HTML/CSS/JS)
├── internal/
│   ├── api/                 # JSON + CSV handlers
│   ├── server/              # mux and routes
│   ├── store/               # in-memory read model
│   ├── config/
│   └── models/
├── pkg/version/             # ldflags-friendly version
├── data/checks.json         # sample controls
└── Dockerfile
```

See [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md) for a concise layer overview.

## Run

Prerequisite: **Go 1.23+** on your `PATH`.

```bash
go run ./cmd/server
```

Open `http://localhost:9090` (same process serves the UI and `/api/v1/*`).

| Variable | Default | Meaning |
|----------|---------|---------|
| `CHECKS_PATH` | `data/checks.json` | Path to controls JSON |
| `PORT` | `9090` | Listen port (`:9090` accepted) |

## Tests

```bash
go test ./...
```

After changing `go.mod` or `go.sum`, run `go mod tidy` with Go 1.23+ on your `PATH`, then run tests again.

## HTTP API (summary)

| Method | Path | Notes |
|--------|------|-------|
| GET | `/health` | `{"status":"ok"}` |
| GET | `/api/v1/version` | Build metadata (`pkg/version`) |
| GET | `/api/v1/checks` | Optional `?tag=CIS` (framework tag) |
| GET | `/api/v1/checks/{id}` | Single row |
| GET | `/api/v1/drift` | Rows where observed ≠ expected |
| GET | `/api/v1/summary` | Counts and per-tag histogram |
| GET | `/api/v1/frameworks` | Sorted `{tag,count}` rows |
| GET | `/api/v1/export/checks.csv` | Optional `?tag=` |

## Docker

```bash
docker build -t compliance .
docker run --rm -p 9090:9090 compliance
```

The image runs as a non-privileged user; sample `data/` is included in the build context.

## License

MIT. See [LICENSE](./LICENSE).
