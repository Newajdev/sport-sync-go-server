#!/usr/bin/env bash
set -euo pipefail

BASE="${BASE_URL:-http://localhost:8080}"
TS=$(date +%s)
PASS=0
FAIL=0

pass() { echo "PASS: $1"; PASS=$((PASS + 1)); }
fail() { echo "FAIL: $1"; FAIL=$((FAIL + 1)); }

assert_eq() {
  if [ "$2" = "$3" ]; then pass "$1"; else fail "$1 (expected '$3', got '$2')"; fi
}

assert_http() {
  if [ "$2" = "$3" ]; then pass "$1"; else fail "$1 (expected HTTP $3, got $2)"; fi
}

assert_contains() {
  if echo "$3" | grep -q "$2"; then pass "$1"; else fail "$1 (missing '$2')"; fi
}

assert_not_contains() {
  if echo "$3" | grep -q "$2"; then fail "$1 (found '$2')"; else pass "$1"; fi
}

echo "=== SpotSync E2E ==="

HEALTH=$(curl -s "$BASE/health")
assert_eq "health" "$HEALTH" "running"

DRIVER_EMAIL="driver${TS}@test.com"
REG=$(curl -s -w "\n%{http_code}" -X POST "$BASE/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Test Driver\",\"email\":\"$DRIVER_EMAIL\",\"password\":\"securePassword123\",\"role\":\"driver\"}")
REG_BODY=$(echo "$REG" | sed '$d')
REG_CODE=$(echo "$REG" | tail -n1)
assert_http "register driver" "$REG_CODE" "201"
assert_contains "register success envelope" '"success":true' "$REG_BODY"
assert_not_contains "register hides password" '"password"' "$REG_BODY"

ADMIN_EMAIL="admin${TS}@test.com"
curl -s -X POST "$BASE/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Test Admin\",\"email\":\"$ADMIN_EMAIL\",\"password\":\"securePassword123\",\"role\":\"admin\"}" > /dev/null

DRIVER_LOGIN=$(curl -s -X POST "$BASE/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$DRIVER_EMAIL\",\"password\":\"securePassword123\"}")
DRIVER_TOKEN=$(echo "$DRIVER_LOGIN" | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')
assert_contains "driver token issued" "$DRIVER_TOKEN" "$DRIVER_LOGIN"

ADMIN_LOGIN=$(curl -s -X POST "$BASE/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"securePassword123\"}")
ADMIN_TOKEN=$(echo "$ADMIN_LOGIN" | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')
assert_contains "admin token issued" "$ADMIN_TOKEN" "$ADMIN_LOGIN"

FORB_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/api/v1/zones" \
  -H "Authorization: Bearer $DRIVER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Should Fail","type":"general","total_capacity":5,"price_per_hour":2.00}')
assert_http "driver cannot create zone" "$FORB_CODE" "403"

ZONE=$(curl -s -w "\n%{http_code}" -X POST "$BASE/api/v1/zones" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"E2E Zone","type":"ev_charging","total_capacity":2,"price_per_hour":5.50}')
ZONE_BODY=$(echo "$ZONE" | sed '$d')
ZONE_CODE=$(echo "$ZONE" | tail -n1)
assert_http "admin create zone" "$ZONE_CODE" "201"
ZONE_ID=$(echo "$ZONE_BODY" | sed -n 's/.*"id":\([0-9]*\).*/\1/p' | head -1)

ZONES_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE/api/v1/zones")
assert_http "public list zones" "$ZONES_CODE" "200"

ZONE_GET=$(curl -s "$BASE/api/v1/zones/$ZONE_ID")
assert_contains "get zone available_spots" '"available_spots":2' "$ZONE_GET"

RES=$(curl -s -w "\n%{http_code}" -X POST "$BASE/api/v1/reservations" \
  -H "Authorization: Bearer $DRIVER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"zone_id\":$ZONE_ID,\"license_plate\":\"ABC-1234\"}")
RES_BODY=$(echo "$RES" | sed '$d')
RES_CODE=$(echo "$RES" | tail -n1)
assert_http "create reservation" "$RES_CODE" "201"
assert_contains "reservation status active" '"status":"active"' "$RES_BODY"
RES_ID=$(echo "$RES_BODY" | sed -n 's/.*"id":\([0-9]*\).*/\1/p' | head -1)

ZONE_AFTER=$(curl -s "$BASE/api/v1/zones/$ZONE_ID")
assert_contains "available_spots decremented" '"available_spots":1' "$ZONE_AFTER"

MY_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE/api/v1/reservations/my-reservations" \
  -H "Authorization: Bearer $DRIVER_TOKEN")
assert_http "my reservations" "$MY_CODE" "200"

ALL_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE/api/v1/reservations" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
assert_http "admin list reservations" "$ALL_CODE" "200"

DALL_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE/api/v1/reservations" \
  -H "Authorization: Bearer $DRIVER_TOKEN")
assert_http "driver blocked from admin list" "$DALL_CODE" "403"

CANCEL_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$BASE/api/v1/reservations/$RES_ID" \
  -H "Authorization: Bearer $DRIVER_TOKEN")
assert_http "cancel reservation" "$CANCEL_CODE" "200"

ZONE_RESTORED=$(curl -s "$BASE/api/v1/zones/$ZONE_ID")
assert_contains "available_spots restored" '"available_spots":2' "$ZONE_RESTORED"

curl -s -X POST "$BASE/api/v1/reservations" -H "Authorization: Bearer $DRIVER_TOKEN" \
  -H "Content-Type: application/json" -d "{\"zone_id\":$ZONE_ID,\"license_plate\":\"FULL-1\"}" > /dev/null
curl -s -X POST "$BASE/api/v1/reservations" -H "Authorization: Bearer $DRIVER_TOKEN" \
  -H "Content-Type: application/json" -d "{\"zone_id\":$ZONE_ID,\"license_plate\":\"FULL-2\"}" > /dev/null
FULL_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/api/v1/reservations" \
  -H "Authorization: Bearer $DRIVER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"zone_id\":$ZONE_ID,\"license_plate\":\"FULL-3\"}")
assert_http "zone full returns 409" "$FULL_CODE" "409"

UNAUTH_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE/api/v1/reservations/my-reservations")
assert_http "missing token returns 401" "$UNAUTH_CODE" "401"

SMALL=$(curl -s -X POST "$BASE/api/v1/zones" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Concurrency Zone","type":"general","total_capacity":1,"price_per_hour":1.00}')
SMALL_ID=$(echo "$SMALL" | sed -n 's/.*"id":\([0-9]*\).*/\1/p' | head -1)

D2_EMAIL="driver2_${TS}@test.com"
curl -s -X POST "$BASE/api/v1/auth/register" -H "Content-Type: application/json" \
  -d "{\"name\":\"Driver Two\",\"email\":\"$D2_EMAIL\",\"password\":\"securePassword123\",\"role\":\"driver\"}" > /dev/null
D2_LOGIN=$(curl -s -X POST "$BASE/api/v1/auth/login" -H "Content-Type: application/json" \
  -d "{\"email\":\"$D2_EMAIL\",\"password\":\"securePassword123\"}")
D2_TOKEN=$(echo "$D2_LOGIN" | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')

TMPDIR="${TMPDIR:-/tmp}"
A_FILE="$TMPDIR/spotsync_a_${TS}.txt"
B_FILE="$TMPDIR/spotsync_b_${TS}.txt"

curl -s -w "\n%{http_code}" -X POST "$BASE/api/v1/reservations" \
  -H "Authorization: Bearer $DRIVER_TOKEN" -H "Content-Type: application/json" \
  -d "{\"zone_id\":$SMALL_ID,\"license_plate\":\"CONC-A\"}" > "$A_FILE" &
curl -s -w "\n%{http_code}" -X POST "$BASE/api/v1/reservations" \
  -H "Authorization: Bearer $D2_TOKEN" -H "Content-Type: application/json" \
  -d "{\"zone_id\":$SMALL_ID,\"license_plate\":\"CONC-B\"}" > "$B_FILE" &
wait

CODE_A=$(tail -n1 "$A_FILE")
CODE_B=$(tail -n1 "$B_FILE")
SUCCESS_COUNT=0
[ "$CODE_A" = "201" ] && SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
[ "$CODE_B" = "201" ] && SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
if [ "$SUCCESS_COUNT" -eq 1 ]; then
  pass "concurrency: exactly one booking succeeds"
else
  fail "concurrency: expected one 201, got $CODE_A and $CODE_B"
fi

echo ""
echo "=============================="
echo "E2E RESULT: $PASS passed, $FAIL failed"
echo "=============================="

[ "$FAIL" -eq 0 ]
