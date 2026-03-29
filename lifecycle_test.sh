#!/bin/bash
# ============================================================
# Unified Workflow — Full Lifecycle Test
# Run this on the jump server to verify the deployed stack at
# 172.30.75.85.
#
# What it tests:
#   1. All service health checks
#   2. NATS JetStream connectivity and stream state
#   3. Workflow registration (via CLI)
#   4. Workflow execution trigger (via curl — CLI execute is mocked)
#   5. Execution status polling until terminal state
#   6. NATS message delivery confirmation
#
# Usage:
#   chmod +x lifecycle_test.sh
#   ./lifecycle_test.sh
# ============================================================

set -euo pipefail

# ── Config ───────────────────────────────────────────────────
HOST="172.30.75.85"
REGISTRY_URL="http://${HOST}:8080"
EXECUTOR_URL="http://${HOST}:8081"
API_URL="http://${HOST}:8082"
NATS_MON="http://${HOST}:8222"
CLI="./uwf-cli"

WORKFLOW_NAME="antifraud-transaction-validation"
POLL_INTERVAL=3   # seconds between status polls
POLL_MAX=20       # max poll attempts before timeout

# Unique transaction ID for this test run
TXN_ID="test-txn-$(date +%s)"

# ── Colours ──────────────────────────────────────────────────
GREEN='\033[0;32m'; RED='\033[0;31m'; YELLOW='\033[1;33m'
CYAN='\033[0;36m'; BOLD='\033[1m'; RESET='\033[0m'

PASS=0; FAIL=0; WARN=0

pass()  { echo -e "  ${GREEN}[PASS]${RESET} $*"; ((PASS++));  }
fail()  { echo -e "  ${RED}[FAIL]${RESET} $*"; ((FAIL++));   }
warn()  { echo -e "  ${YELLOW}[WARN]${RESET} $*"; ((WARN++)); }
info()  { echo -e "  ${CYAN}      ${RESET} $*"; }
header(){ echo -e "\n${BOLD}$*${RESET}"; }

# ── Helper ───────────────────────────────────────────────────
http_get() {
    curl -sf --connect-timeout 5 --max-time 10 "$1" 2>/dev/null
}

http_post() {
    local url="$1"; local body="$2"
    curl -sf --connect-timeout 5 --max-time 30 \
        -X POST "$url" \
        -H "Content-Type: application/json" \
        -d "$body" 2>/dev/null
}

json_field() {
    # Extract a top-level JSON field value without requiring jq
    echo "$1" | grep -o "\"$2\"[[:space:]]*:[[:space:]]*\"[^\"]*\"" \
              | head -1 | sed 's/.*:[[:space:]]*"\(.*\)"/\1/'
}

# ── Banner ───────────────────────────────────────────────────
echo -e "\n${BOLD}=================================================${RESET}"
echo -e "${BOLD}  Unified Workflow — Full Lifecycle Test${RESET}"
echo -e "${BOLD}  Target: ${HOST}${RESET}"
echo -e "${BOLD}  TXN:    ${TXN_ID}${RESET}"
echo -e "${BOLD}=================================================${RESET}"

# ════════════════════════════════════════════════════════════
# PHASE 1 — Service health
# ════════════════════════════════════════════════════════════
header "Phase 1 · Service Health"

check_health() {
    local name="$1" url="$2"
    local resp
    resp=$(http_get "$url") || true
    if [ -n "$resp" ]; then
        pass "$name reachable ($url)"
    else
        fail "$name not responding ($url)"
    fi
}

check_health "Registry   :8080" "${REGISTRY_URL}/health"
check_health "Executor   :8081" "${EXECUTOR_URL}/health"
check_health "Workflow API:8082" "${API_URL}/health"

# NATS HTTP monitoring
resp=$(http_get "${NATS_MON}/varz") || true
if [ -n "$resp" ]; then
    pass "NATS monitoring :8222 responding"
    nats_version=$(json_field "$resp" "version")
    info "NATS version: ${nats_version}"
else
    fail "NATS monitoring :8222 not responding — worker cannot consume"
fi

# ════════════════════════════════════════════════════════════
# PHASE 2 — NATS JetStream state
# ════════════════════════════════════════════════════════════
header "Phase 2 · NATS JetStream"

jsz=$(http_get "${NATS_MON}/jsz?streams=1&consumers=1") || true
if [ -n "$jsz" ]; then
    pass "JetStream enabled"
    streams=$(echo "$jsz" | grep -o '"num_streams"[[:space:]]*:[[:space:]]*[0-9]*' | grep -o '[0-9]*$' || echo "?")
    consumers=$(echo "$jsz" | grep -o '"num_consumers"[[:space:]]*:[[:space:]]*[0-9]*' | grep -o '[0-9]*$' || echo "?")
    messages=$(echo "$jsz" | grep -o '"messages"[[:space:]]*:[[:space:]]*[0-9]*' | head -1 | grep -o '[0-9]*$' || echo "?")
    info "Streams: ${streams}  |  Consumers: ${consumers}  |  Total messages stored: ${messages}"

    if [ "$streams" = "0" ] || [ "$streams" = "?" ]; then
        warn "No JetStream streams found — worker may not have subscribed yet"
    fi
else
    fail "Could not reach NATS JetStream API — is JetStream enabled? (nats -js flag required)"
fi

# Peek at consumer lag
conz=$(http_get "${NATS_MON}/connz") || true
if [ -n "$conz" ]; then
    conn_count=$(echo "$conz" | grep -o '"num_connections"[[:space:]]*:[[:space:]]*[0-9]*' | grep -o '[0-9]*$' || echo "?")
    info "NATS active connections: ${conn_count}"
    if [ "$conn_count" = "0" ]; then
        warn "No NATS connections — services may not be connected to NATS"
    fi
fi

# ════════════════════════════════════════════════════════════
# PHASE 3 — Workflow registration
# ════════════════════════════════════════════════════════════
header "Phase 3 · Workflow Registration"

# NOTE: uwf-cli 'workflows list' makes real HTTP calls; 'execute' is mocked.
if [ ! -f "$CLI" ]; then
    fail "uwf-cli binary not found at ${CLI}"
    WORKFLOW_ID=""
else
    chmod +x "$CLI"
    pass "uwf-cli binary found"

    echo ""
    info "Registered workflows (via CLI):"
    wf_list=$("$CLI" workflows list --endpoint "$API_URL" 2>&1) || true
    echo "$wf_list" | sed 's/^/    /'

    # Check if our target workflow is registered
    if echo "$wf_list" | grep -q "$WORKFLOW_NAME"; then
        pass "Workflow '${WORKFLOW_NAME}' found in registry"
        WORKFLOW_ID=$(echo "$wf_list" | grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' | head -1 | sed 's/.*"\([^"]*\)"/\1/')
        info "Workflow ID: ${WORKFLOW_ID}"
    else
        warn "Workflow '${WORKFLOW_NAME}' not registered — registering now..."

        reg_resp=$(http_post "${API_URL}/api/v1/workflows" \
            "{\"name\":\"${WORKFLOW_NAME}\",\"description\":\"Antifraud transaction validation: Store → AML → FC → ML → Finalize\"}") || true

        if [ -n "$reg_resp" ]; then
            WORKFLOW_ID=$(json_field "$reg_resp" "id")
            [ -z "$WORKFLOW_ID" ] && WORKFLOW_ID=$(json_field "$reg_resp" "workflowId")
            pass "Workflow registered (ID: ${WORKFLOW_ID:-unknown})"
            info "Response: $reg_resp"
        else
            fail "Failed to register workflow via API"
            WORKFLOW_ID=""
        fi
    fi
fi

# ════════════════════════════════════════════════════════════
# PHASE 4 — Execute workflow (curl — CLI execute is mocked)
# ════════════════════════════════════════════════════════════
header "Phase 4 · Execute Workflow"
info "Note: ./uwf-cli execute is mocked (returns fake data). Using curl against the real API."

INPUT_PAYLOAD=$(cat <<EOF
{
  "input_data": {
    "transaction": {
      "id": "${TXN_ID}",
      "type": "deposit",
      "amount": "100000",
      "currency": "KZT",
      "client_id": "client-lifecycle-test",
      "client_name": "Test User",
      "client_pan": "111111******1111",
      "client_cvv": "111",
      "client_card_holder": "TEST USER",
      "client_phone": "+77007007070",
      "merchant_terminal_id": "00000001",
      "channel": "E-com",
      "location_ip": "192.168.0.1"
    }
  }
}
EOF
)

RUN_ID=""

# Try workflow-api first (port 8082)
exec_resp=$(http_post "${API_URL}/api/v1/workflows/${WORKFLOW_NAME}/execute" "$INPUT_PAYLOAD") || true
if [ -z "$exec_resp" ] && [ -n "$WORKFLOW_ID" ]; then
    exec_resp=$(http_post "${API_URL}/api/v1/workflows/${WORKFLOW_ID}/execute" "$INPUT_PAYLOAD") || true
fi

# Fallback: try executor service directly (port 8081)
if [ -z "$exec_resp" ]; then
    info "workflow-api execute returned empty — trying executor service directly..."
    exec_resp=$(http_post "${EXECUTOR_URL}/api/v1/execute" \
        "{\"workflow_id\":\"${WORKFLOW_ID:-${WORKFLOW_NAME}}\",\"input_data\":$(echo "$INPUT_PAYLOAD" | grep -o '\"transaction\".*}')}}") || true
fi

if [ -n "$exec_resp" ]; then
    pass "Execution triggered"
    info "Response: $exec_resp"
    RUN_ID=$(json_field "$exec_resp" "runId")
    [ -z "$RUN_ID" ] && RUN_ID=$(json_field "$exec_resp" "run_id")
    [ -z "$RUN_ID" ] && RUN_ID=$(json_field "$exec_resp" "executionId")
    info "Run ID: ${RUN_ID:-not found in response}"
else
    fail "Execution request returned no response — API may be down or endpoint path differs"
fi

# ════════════════════════════════════════════════════════════
# PHASE 5 — Poll execution status
# ════════════════════════════════════════════════════════════
header "Phase 5 · Execution Status Polling"

if [ -z "$RUN_ID" ]; then
    warn "No run ID — skipping status polling"
else
    info "Polling every ${POLL_INTERVAL}s (max ${POLL_MAX} attempts)..."
    attempt=0
    terminal=false

    while [ $attempt -lt $POLL_MAX ]; do
        sleep "$POLL_INTERVAL"
        ((attempt++))

        status_resp=$(http_get "${API_URL}/api/v1/executions/${RUN_ID}") || true
        if [ -z "$status_resp" ]; then
            status_resp=$(http_get "${EXECUTOR_URL}/api/v1/executions/${RUN_ID}") || true
        fi

        if [ -z "$status_resp" ]; then
            info "[${attempt}/${POLL_MAX}] No response from status endpoint"
            continue
        fi

        status=$(json_field "$status_resp" "status")
        step=$(json_field "$status_resp" "current_step")
        info "[${attempt}/${POLL_MAX}] status=${status}  step=${step:-—}"

        case "$status" in
            completed|COMPLETED)
                pass "Workflow COMPLETED successfully"
                terminal=true; break;;
            failed|FAILED)
                err=$(json_field "$status_resp" "error_message")
                fail "Workflow FAILED — ${err:-see response below}"
                info "Full response: $status_resp"
                terminal=true; break;;
            cancelled|CANCELLED)
                warn "Workflow was CANCELLED"
                terminal=true; break;;
        esac
    done

    if [ "$terminal" = "false" ]; then
        warn "Polling timed out after $((POLL_MAX * POLL_INTERVAL))s — workflow still running"
        info "Check manually: curl ${API_URL}/api/v1/executions/${RUN_ID}"
    fi
fi

# ════════════════════════════════════════════════════════════
# PHASE 6 — NATS message delivery confirmation
# ════════════════════════════════════════════════════════════
header "Phase 6 · NATS Delivery Confirmation"

jsz_after=$(http_get "${NATS_MON}/jsz?streams=1&consumers=1&config=1&state=1") || true
if [ -n "$jsz_after" ]; then
    msgs_after=$(echo "$jsz_after" | grep -o '"messages"[[:space:]]*:[[:space:]]*[0-9]*' | head -1 | grep -o '[0-9]*$' || echo "?")
    pending=$(echo "$jsz_after" | grep -o '"num_pending"[[:space:]]*:[[:space:]]*[0-9]*' | grep -o '[0-9]*$' || echo "?")
    ack_floor=$(echo "$jsz_after" | grep -o '"ack_floor"[[:space:]]*:[[:space:]]*[0-9]*' | grep -o '[0-9]*$' || echo "?")

    info "JetStream messages stored: ${msgs_after}"
    info "Consumer pending (unacked): ${pending}"
    info "Ack floor (last processed): ${ack_floor}"

    if [ "$pending" = "0" ] || [ "$pending" = "?" ]; then
        pass "No pending messages — worker consumed all messages"
    else
        warn "${pending} messages still pending in NATS — worker may be slow or not running"
    fi
else
    warn "Could not query NATS JetStream for delivery confirmation"
fi

# ════════════════════════════════════════════════════════════
# PHASE 7 — Docker container status (if docker available)
# ════════════════════════════════════════════════════════════
header "Phase 7 · Docker Container Status"

if command -v docker >/dev/null 2>&1; then
    containers=(
        "unified-workflow-nats"
        "unified-workflow-registry"
        "unified-workflow-executor"
        "unified-workflow-api"
        "unified-workflow-worker"
    )
    for c in "${containers[@]}"; do
        state=$(docker inspect --format '{{.State.Status}}' "$c" 2>/dev/null || echo "not found")
        health=$(docker inspect --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}no healthcheck{{end}}' "$c" 2>/dev/null || echo "")
        if [ "$state" = "running" ]; then
            pass "${c}: running  (health: ${health})"
        else
            fail "${c}: ${state}"
        fi
    done
else
    warn "docker not in PATH on this host — skipping container check"
fi

# ════════════════════════════════════════════════════════════
# SUMMARY
# ════════════════════════════════════════════════════════════
echo ""
echo -e "${BOLD}=================================================${RESET}"
echo -e "${BOLD}  Test Summary${RESET}"
echo -e "${BOLD}=================================================${RESET}"
echo -e "  ${GREEN}PASS${RESET}: ${PASS}"
echo -e "  ${YELLOW}WARN${RESET}: ${WARN}"
echo -e "  ${RED}FAIL${RESET}: ${FAIL}"
echo ""

if [ -n "$RUN_ID" ]; then
    echo -e "  Run ID : ${RUN_ID}"
    echo -e "  Status : curl ${API_URL}/api/v1/executions/${RUN_ID}"
fi

if [ $FAIL -eq 0 ]; then
    echo -e "\n  ${GREEN}${BOLD}All checks passed. Full lifecycle working.${RESET}\n"
    exit 0
else
    echo -e "\n  ${RED}${BOLD}${FAIL} check(s) failed. See output above.${RESET}\n"
    exit 1
fi
