#!/usr/bin/env bash
# Real-Neo4j HTTPS smoke test for `neo4j-cli query --insecure`.
#
# Boots a Neo4j 5 docker container with HTTPS enabled on port 7473 using a
# freshly-generated self-signed certificate, then exercises the CLI binary
# against it. Asserts the positive (`--insecure` succeeds) and negative
# (default TLS verification fails on a self-signed cert) paths.
#
# Run directly:    bash scripts/test-https.sh
# Or from Go:      NEO4J_HTTPS_TEST=1 go test ./neo4j-cli/query/ -run TestHTTPS_Smoke -v
#
# Skip the build:  SKIP_BUILD=1 bash scripts/test-https.sh   (uses bin/neo4j-cli)

set -euo pipefail

CONTAINER_NAME="neo4j-https-test"
NEO4J_IMAGE="${NEO4J_IMAGE:-neo4j:5}"
NEO4J_PASSWORD="testtest"
HTTPS_PORT=7473
HTTP_PORT=7474
READY_TIMEOUT=60

# Resolve repo root from script location so the script works from any cwd.
SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)
REPO_ROOT=$(cd -- "${SCRIPT_DIR}/.." &> /dev/null && pwd)

tempdir=""
cleanup() {
  docker rm -f "${CONTAINER_NAME}" >/dev/null 2>&1 || true
  if [ -n "${tempdir}" ] && [ -d "${tempdir}" ]; then
    rm -rf "${tempdir}"
  fi
}
trap cleanup EXIT

# --- pre-flight ---------------------------------------------------------------
need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "ERROR: \`$1\` is required but not on PATH. Install it and retry." >&2
    exit 2
  fi
}
need docker
need openssl
need curl

port_in_use() {
  # macOS + Linux compatible: use lsof if available, fall back to nc.
  if command -v lsof >/dev/null 2>&1; then
    lsof -nP -iTCP:"$1" -sTCP:LISTEN >/dev/null 2>&1
  elif command -v nc >/dev/null 2>&1; then
    nc -z localhost "$1" >/dev/null 2>&1
  else
    return 1
  fi
}

for port in "${HTTPS_PORT}" "${HTTP_PORT}"; do
  if port_in_use "${port}"; then
    echo "ERROR: port ${port} is already in use. Stop the listener and retry." >&2
    exit 2
  fi
done

# Make sure we don't collide with a leftover container from a previous run.
docker rm -f "${CONTAINER_NAME}" >/dev/null 2>&1 || true

# --- build binary -------------------------------------------------------------
BIN="${REPO_ROOT}/bin/neo4j-cli"
if [ "${SKIP_BUILD:-0}" != "1" ]; then
  echo "==> building neo4j-cli"
  (cd "${REPO_ROOT}" && make build >/dev/null)
fi
if [ ! -x "${BIN}" ]; then
  echo "ERROR: ${BIN} not found or not executable. Run \`make build\` or unset SKIP_BUILD." >&2
  exit 2
fi

# --- self-signed cert ---------------------------------------------------------
tempdir=$(mktemp -d)
echo "==> generating self-signed cert under ${tempdir}"
openssl req -x509 -newkey rsa:2048 -nodes \
  -keyout "${tempdir}/private.key" \
  -out "${tempdir}/public.crt" \
  -days 30 \
  -subj "/CN=localhost" \
  >/dev/null 2>&1

mkdir -p "${tempdir}/trusted" "${tempdir}/revoked"

# Neo4j container runs as a non-root user (uid 7474). Make all cert files
# world-readable so the in-container reader can pick them up; the private key
# stays at 0644 because the bind mount is read-only and the host dir is
# already inside a private mktemp tree.
chmod 0644 "${tempdir}/private.key" "${tempdir}/public.crt"
chmod 0755 "${tempdir}" "${tempdir}/trusted" "${tempdir}/revoked"

# --- boot Neo4j ---------------------------------------------------------------
echo "==> booting ${NEO4J_IMAGE} with HTTPS enabled"
docker run -d --rm \
  --name "${CONTAINER_NAME}" \
  -p "${HTTPS_PORT}":7473 \
  -p "${HTTP_PORT}":7474 \
  -e NEO4J_AUTH="neo4j/${NEO4J_PASSWORD}" \
  -e NEO4J_server_https_enabled=true \
  -e NEO4J_dbms_ssl_policy_https_enabled=true \
  -e NEO4J_dbms_ssl_policy_https_base__directory=/ssl \
  -e NEO4J_dbms_ssl_policy_https_private__key=private.key \
  -e NEO4J_dbms_ssl_policy_https_public__certificate=public.crt \
  -e NEO4J_dbms_ssl_policy_https_client__auth=NONE \
  -v "${tempdir}":/ssl:ro \
  "${NEO4J_IMAGE}" >/dev/null

echo "==> waiting up to ${READY_TIMEOUT}s for HTTPS readiness on :${HTTPS_PORT}"
deadline=$(($(date +%s) + READY_TIMEOUT))
ready=0
while [ "$(date +%s)" -lt "${deadline}" ]; do
  if curl -sf -k "https://localhost:${HTTPS_PORT}" >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 2
done
if [ "${ready}" -ne 1 ]; then
  echo "ERROR: Neo4j HTTPS endpoint did not become ready within ${READY_TIMEOUT}s" >&2
  echo "--- last 50 lines of container logs ---" >&2
  docker logs --tail 50 "${CONTAINER_NAME}" >&2 || true
  exit 1
fi

# --- positive: --insecure succeeds --------------------------------------------
echo "==> [positive] expecting --insecure to succeed"
positive_out=$("${BIN}" query \
  --uri "https://localhost:${HTTPS_PORT}" \
  --insecure \
  -u neo4j \
  -p "${NEO4J_PASSWORD}" \
  "RETURN 1 AS n")
if ! grep -q '1' <<< "${positive_out}"; then
  echo "ERROR: positive output did not contain '1'. Got:" >&2
  echo "${positive_out}" >&2
  exit 1
fi
echo "    OK — stdout contained '1'"

# --- negative: default verification rejects self-signed -----------------------
echo "==> [negative] expecting default (no --insecure) to fail TLS verification"
set +e
negative_stderr=$("${BIN}" query \
  --uri "https://localhost:${HTTPS_PORT}" \
  -u neo4j \
  -p "${NEO4J_PASSWORD}" \
  "RETURN 1" 2>&1 >/dev/null)
negative_rc=$?
set -e

if [ "${negative_rc}" -eq 0 ]; then
  echo "ERROR: negative invocation unexpectedly succeeded (rc=0). stderr:" >&2
  echo "${negative_stderr}" >&2
  exit 1
fi

# Match TLS / x509 / certificate (case-insensitive) — wording varies by Go version.
if ! grep -qiE 'tls|x509|certificate' <<< "${negative_stderr}"; then
  echo "ERROR: negative stderr did not mention TLS/x509/certificate. Got:" >&2
  echo "${negative_stderr}" >&2
  exit 1
fi
echo "    OK — exited non-zero with TLS-related stderr"

echo "OK: --insecure verified end-to-end"
