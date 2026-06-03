#!/bin/sh
# run application from package bin directory
cd bin || exit 1
exec ./v4poolkb \
  -port "${V4_POOL_KB_PORT:-8087}" \
  -jwtkey "$V4_POOL_KB_JWT_KEY" \
  -region "${V4_POOL_KB_AWS_REGION:-us-east-1}" \
  -kbid "$V4_POOL_KB_ID" \
  -threshold "${V4_POOL_KB_SCORE_THRESHOLD:-0}" \
  -limit "${V4_POOL_KB_DEFAULT_LIMIT:-20}" \
  -detailresource "${V4_KB_DETAIL_BASE}"
