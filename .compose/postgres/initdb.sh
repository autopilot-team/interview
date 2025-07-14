#!/usr/bin/env bash 

set -euo pipefail

# scripts for initializing db, mimick commands in Taskfile before but on init script postgres.
# this script takes inputs from envs from Taskfile.

# parse from envs "(string)* -> array of string"
read -ra DB_NO_LIVE_TEST <<< "$DATABASES_WITHOUT_LIVE_TEST"
read -ra DB_WITH_LIVE_TEST <<< "$DATABASES_WITH_LIVE_TEST"

for DB in "${DB_NO_LIVE_TEST[@]}"; do
  echo "Creating database: $DB"
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
    -c "CREATE DATABASE \"$DB\";"
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
    -c "GRANT ALL PRIVILEGES ON DATABASE \"$DB\" TO \"$POSTGRES_USER\";"
done

for DB in "${DB_WITH_LIVE_TEST[@]}"; do
  for SUFFIX in live test; do
    DB_NAME="${DB}_${SUFFIX}"
    echo "Creating database: $DB_NAME"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
      -c "CREATE DATABASE \"$DB_NAME\";"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
      -c "GRANT ALL PRIVILEGES ON DATABASE \"$DB_NAME\" TO \"$POSTGRES_USER\";"
  done
done
