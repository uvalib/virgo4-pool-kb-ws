# virgo4-pool-kb-ws

Virgo 4 pool web service that searches the UVA digital image Bedrock Knowledge Base and returns `v4api.PoolResult` responses for `virgo4-search-ws`.

## Endpoints

- `GET /`, `GET /version` — service version
- `GET /healthcheck` — dependency health
- `GET /identify` — pool identity for search-ws
- `GET /api/providers` — access URL providers (empty for KB pool)
- `POST /api/search` — `v4api.SearchRequest` → `v4api.PoolResult` (JWT required)
- `POST /api/search/facets` — empty facet list (JWT required)
- `GET /api/resource/:id` — proxies to Solr images pool `/api/resource/{id}` (JWT required)
- `GET /uva-lib:...` — same proxy for `pool.url` + identifier requests

## Configuration

Precedence: **defaults < environment < command-line flags**.

| Environment variable | Description |
| --- | --- |
| `V4_POOL_KB_CONFIG` | Path to YAML config file |
| `V4_POOL_KB_PORT` | HTTP port (default `8087`) |
| `V4_POOL_KB_JWT_KEY` | JWT HMAC secret (required) |
| `V4_POOL_KB_AWS_REGION` | AWS region (default `us-east-1`) |
| `V4_POOL_KB_ID` | Bedrock knowledge base id |
| `V4_POOL_KB_SCORE_THRESHOLD` | Minimum retrieval score |
| `V4_POOL_KB_DEFAULT_LIMIT` | Default retrieval result count |
| `V4_KB_DETAIL_BASE` | Solr images pool `/api/resource/` URL base; `{id}` is appended; KB pool proxies the response with the caller JWT |

## Local run

```bash
cp scripts/dev.local.sh.example scripts/dev.local.sh
./scripts/dev.local.sh
```

In virgo's db, add `virgo4.sources` as `images_kb` with `http://localhost:8087`.
Use the same `V4_JWT_KEY` as `virgo4-search-ws` / `virgo4-client`.

VS Code tasks include server startup for **KB pool (bedrock)** and **KB pool (mock)** plus build and test tasks (`.vscode/tasks.json`).

```bash
export V4_POOL_KB_JWT_KEY=jwt-key
./package/scripts/entry.sh
```

## Tests

```bash
go test ./...
```
