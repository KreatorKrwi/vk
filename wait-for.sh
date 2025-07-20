#!/bin/sh

set -e

host="$1"
port="$2"
shift 2
cmd="$@"

until nc -z "$host" "$port"; do
  echo "Waiting for database $host:$port..."
  sleep 1
done

until pg_isready -h "$host" -p "$port" -U "$DB_USER" -d "$DB_NAME"; do
  echo "PostgreSQL is not ready yet..."
  sleep 1
done

if [ -d "/app/migrations" ]; then
  echo "Applying migrations..."
  for migration in /app/migrations/*.up.sql; do
    echo "Applying $(basename "$migration")"
    export PGPASSWORD="$DB_PASSWORD"
    psql -h "$host" -p "$port" -U "$DB_USER" -d "$DB_NAME" -f "$migration"
  done
fi

exec $cmd