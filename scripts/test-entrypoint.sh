#!/bin/sh
set -eu

apk add --no-cache curl jq >/dev/null
sh ./scripts/wait-for-api.sh
sh ./scripts/ci-flows.sh
