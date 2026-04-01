#!/bin/sh
set -eu

BASE_URL="${BASE_URL:-http://api:8080}"
MAX_ATTEMPTS="${MAX_ATTEMPTS:-60}"
SLEEP_SECONDS="${SLEEP_SECONDS:-2}"

attempt=1
while [ "$attempt" -le "$MAX_ATTEMPTS" ]; do
	if curl -fsS "$BASE_URL/healthz" >/dev/null; then
		echo "API is ready"
		exit 0
	fi

	echo "Waiting for API at $BASE_URL/healthz (attempt $attempt/$MAX_ATTEMPTS)"
	attempt=$((attempt + 1))
	sleep "$SLEEP_SECONDS"
done

echo "API did not become ready in time" >&2
exit 1
