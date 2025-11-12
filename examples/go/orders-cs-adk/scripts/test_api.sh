#!/usr/bin/env bash
set -eu
SERVER="${SERVER:-http://localhost:8080}"
curl -s -w "\n%{http_code}\n" "$SERVER/health"
curl -s -w "\n%{http_code}\n" -H 'Content-Type: application/json' -d '{"query":"查询订单 20251112002 物流"}' "$SERVER/chat"
curl -s -w "\n%{http_code}\n" -H 'Content-Type: application/json' -d '{"id":"demo-verify","order_id":"20251112001","title":"某公司","tax_id":"123456"}' "$SERVER/approval/invoice/start"
curl -s -w "\n%{http_code}\n" -H 'Content-Type: application/json' -d '{"id":"demo-verify","order_id":"20251112001","title":"某公司","tax_id":"123456"}' "$SERVER/approval/invoice/resume"
