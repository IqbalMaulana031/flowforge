#!/usr/bin/env sh
set -eu

direction="${1:-up}"

if ! command -v migrate >/dev/null 2>&1; then
  echo "golang-migrate CLI is not installed. Install it before running migrations." >&2
  exit 1
fi

base_database_url="postgres://${POSTGRES_USER:-flowforge}:${POSTGRES_PASSWORD:-secret}@${POSTGRES_HOST:-localhost}:${POSTGRES_PORT:-5432}/${POSTGRES_NAME:-flowforge_db}?sslmode=${POSTGRES_SSLMODE:-disable}"

run_migration() {
  path="$1"
  table="$2"
  command="$3"
  migrate -path "$path" -database "${base_database_url}&x-migrations-table=${table}" "$command"
}

run_up() {
  run_migration db/schemas schema_migrations_schemas up
  run_migration db/migrations/extensions schema_migrations_extensions up
  run_migration db/migrations/master schema_migrations_master up
  run_migration db/migrations/auth schema_migrations_auth up
  run_migration db/migrations/workflow schema_migrations_workflow up
  run_migration db/migrations/execution schema_migrations_execution up
  run_migration db/migrations/scheduler schema_migrations_scheduler up
}

run_down() {
  run_migration db/migrations/scheduler schema_migrations_scheduler down
  run_migration db/migrations/execution schema_migrations_execution down
  run_migration db/migrations/workflow schema_migrations_workflow down
  run_migration db/migrations/auth schema_migrations_auth down
  run_migration db/migrations/master schema_migrations_master down
  run_migration db/migrations/extensions schema_migrations_extensions down
  run_migration db/schemas schema_migrations_schemas down
}

case "$direction" in
  up)
    run_up
    ;;
  down)
    run_down
    ;;
  *)
    echo "Unsupported migration direction: $direction. Use 'up' or 'down'." >&2
    exit 1
    ;;
esac
