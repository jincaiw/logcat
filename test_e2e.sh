#!/bin/bash
set -e
cd /Volumes/ex-data/jason.wa/codebase/golog
pkill -f "./logcat" 2>/dev/null || true
sleep 1
rm -f data/logcat.db

LOGCAT_ADMIN_PASSWORD="Test123456" ./logcat --config configs/config.yaml > /tmp/logcat_e2e3.log 2>&1 &
SERVER_PID=$!
sleep 4

echo "=== 1. Health Check ==="
curl -s http://localhost:8080/healthz
echo ""

echo "=== 2. Ready Check ==="
curl -s http://localhost:8080/readyz
echo ""

echo "=== 3. Login ==="
LOGIN_RESP=$(curl -s -c /tmp/cookies.txt -X POST http://localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"Test123456"}')
echo "$LOGIN_RESP"

echo ""
echo "=== 4. Get Current User ==="
curl -s -b /tmp/cookies.txt http://localhost:8080/api/auth/me
echo ""

echo "=== 5. Create Device ==="
curl -s -b /tmp/cookies.txt -X POST http://localhost:8080/api/devices \
  -H 'Content-Type: application/json' \
  -d '{"name":"test_firewall","ip_address":"192.168.1.100","enabled":true}'
echo ""

echo "=== 6. List Devices ==="
curl -s -b /tmp/cookies.txt "http://localhost:8080/api/devices?page=1&pageSize=10"
echo ""

echo "=== 7. Send Syslog via UDP ==="
for i in 1 2 3; do
  MSG="<134>May 27 23:55:0${i} 192.168.1.100 Firewall: DROP IN=eth0 SRC=10.0.0.${i}"
  echo "$MSG" | nc -u -w1 127.0.0.1 5140
  echo "Sent msg ${i}"
done
sleep 3

echo "=== 8. Query Logs ==="
curl -s -b /tmp/cookies.txt "http://localhost:8080/api/logs?page=1&pageSize=5"
echo ""

echo "=== 9. Dashboard ==="
curl -s -b /tmp/cookies.txt http://localhost:8080/api/dashboard
echo ""

echo "=== 10. Runtime Metrics ==="
curl -s http://localhost:8080/api/metrics/runtime
echo ""

echo "=== 11. Syslog Status ==="
curl -s -b /tmp/cookies.txt http://localhost:8080/api/system/status
echo ""

echo "=== 12. Frontend Index ==="
curl -s -o /dev/null -w "HTTP %{http_code}, size: %{size_download} bytes" http://localhost:8080/
echo ""

kill $SERVER_PID 2>/dev/null
wait $SERVER_PID 2>/dev/null
echo "=== All tests complete ==="