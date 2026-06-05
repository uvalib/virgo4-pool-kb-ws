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
- `GET /{namespace:id}` — same proxy for `pool.url` + identifier requests (e.g. `uva-lib:123`, `tsb:59492`)

## Configuration

Precedence: **defaults < environment < command-line flags**.

| Environment variable | Description |
| --- | --- |
| `V4_POOL_KB_PORT` | HTTP port (default `8087`) |
| `JWT_KEY` | JWT secret (required) |
| `AWS_REGION` | AWS region (default `us-east-1`) |
| `AWS_KB_ID` | Bedrock knowledge base id |
| `AWS_KB_SCORE_THRESHOLD` | Minimum retrieval score |
| `AWS_KB_DEFAULT_LIMIT` | Default retrieval result count |
| `AWS_KB_PROVIDER` | Knowledge base backend: `bedrock` or `mock` (default `bedrock`) |
| `IIIF_BASE_URL` | IIIF image API base. Used with kb `iiif_id` to make image URLs (default `https://iiif.lib.virginia.edu/iiif`) |
| `VIRGO_DETAIL_SOURCE` | Solr images pool `/api/resource/` URL base; `{id}` is appended; KB pool proxies the response with the caller JWT |

## Local run

```bash
cp scripts/dev.local.sh.example scripts/dev.local.sh
./scripts/dev.local.sh
```

In virgo's db, add `virgo4.sources` as `images_kb` with `http://localhost:8087`.
Use the same `V4_JWT_KEY` as `virgo4-search-ws` / `virgo4-client`.

VS Code tasks include server startup for **KB pool (bedrock)** and **KB pool (mock)** plus build and test tasks (`.vscode/tasks.json`).

```bash
export JWT_KEY=jwt-key
./package/scripts/entry.sh
```

## Deployment

Staging ECS and CodePipeline live in `terraform-infrastructure`:

- `virgo4.lib.virginia.edu/ecs-tasks/staging/pool-kb-ws`
- `virgo4.lib.virginia.edu/pipelines/pool-kb-ws`

## Tests

```bash
go test ./...
```
