#!/bin/sh

set -e
MIGRATION_DIR=/migrations
DB_DSN="e_university"
if [ "$1" = "--dryrun" ]; then
  goose -allow-missing -dir ${MIGRATION_DIR} postgres "${DB_DSN}" status
else
  goose -allow-missing -dir ${MIGRATION_DIR} postgres "${DB_DSN}" up
fi