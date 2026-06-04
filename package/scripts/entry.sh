#!/bin/sh
# run application from package bin directory
cd bin || exit 1
exec ./v4poolkb \
  -port "${V4_POOL_KB_PORT:-8087}" \
  -jwtkey "$JWT_KEY" \
  -region "${AWS_REGION:-us-east-1}" \
  -kbid "$AWS_KB_ID" \
  -threshold "${AWS_KB_SCORE_THRESHOLD:-0}" \
  -limit "${AWS_KB_DEFAULT_LIMIT:-20}" \
  -provider "${AWS_KB_PROVIDER:-bedrock}" \
  -iiifbase "${IIIF_BASE_URL}" \
  -detailresource "${VIRGO_DETAIL_SOURCE}"
