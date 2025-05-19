#!/usr/bin/env bash

set -euo pipefail

# scripts for waiting for db, rather than use pg_isready (which is only checks for connectivity)
# checks actual db is exists or not. it fallbacks to pg_isready if both either envs not exists

wait_for_server() {

  echo "waiting for db $POSTGRES_HOST:$POSTGRES_PORT"

  until pg_isready -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" >/dev/null 2>&1; do
    echo "server not ready"
    sleep 2
  done

  echo "db is ready"
}

wait_for_db() {
  
  local db_name=$1

  echo "waiting for db: $db_name"

  until psql -U "$POSTGRES_USER" -d "$db_name" -h "${POSTGRES_HOST:-localhost}" -c '\q' &> /dev/null; do
    echo "$db_name not ready"
    sleep 2
  done

  echo "$db_name is ready"
}

if [[ -z "${DATABASES_WITHOUT_LIVE_TEST:-}" && -z "${DATABASES_WITH_LIVE_TEST:-}" ]]; then
  wait_for_server
  exit 0
fi

# parse from envs "(string)* -> array of string"
read -ra DB_NO_LIVE_TEST <<< "$DATABASES_WITHOUT_LIVE_TEST"
read -ra DB_WITH_LIVE_TEST <<< "$DATABASES_WITH_LIVE_TEST"

# Wait for non live test DBs
if [[ ${#DB_NO_LIVE_TEST[@]} -gt 0 ]]; then
  for DB in "${DB_NO_LIVE_TEST[@]}"; do
    wait_for_db "$DB"
  done
else
  wait_for_server
fi

# Wait for live/test DBs
if [[ ${#DB_WITH_LIVE_TEST[@]} -gt 0 ]]; then
  for DB in "${DB_WITH_LIVE_TEST[@]}"; do
    for SUFFIX in live test; do
      DB_NAME="${DB}_${SUFFIX}"
      wait_for_db "$DB_NAME"
    done
  done
else
  wait_for_server
fi
