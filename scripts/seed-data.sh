#!/usr/bin/env bash
# quick and dirty script to seed some fake transactions for local testing
# not used anywhere in CI, just for manual dev use

set -e

BIN=./moneyback
DB=${1:-./dev.json}

if [ ! -f "$BIN" ]; then
  echo "build the CLI first: go build ./cmd/moneyback"
  exit 1
fi

$BIN -db "$DB" add -account checking -amount 2000 -category salary -desc "paycheck"
$BIN -db "$DB" add -account checking -amount -45.50 -category groceries -desc "trader joes"
$BIN -db "$DB" add -account checking -amount -12.00 -category transit -desc "train pass"
$BIN -db "$DB" add -account savings -amount 500 -category transfer -desc "monthly savings"

echo "seeded $DB"
