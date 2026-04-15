# Compliance viewer — layout

This service follows the same rough shape as larger Go tools: a thin `cmd/` entrypoint, most logic under `internal/`, and link-time metadata in `pkg/version`.

| Path | Role |
|------|------|
| `cmd/server` | `main`, embeds `static/` UI assets |
| `internal/api` | JSON + CSV HTTP handlers |
| `internal/config` | Environment defaults (`CHECKS_PATH`, `PORT`) |
| `internal/server` | Route table (`net/http` patterns; stdlib feature in **Go 1.22+**; this module targets **Go 1.23** in `go.mod`) |
| `internal/store` | In-memory read model over `data/checks.json` |
| `internal/models` | Shared structs |
| `pkg/version` | `-ldflags -X …Version=` at release build time |

The UI is static HTML/CSS/JS served from the embedded filesystem; the API is plain REST under `/api/v1/*`.
