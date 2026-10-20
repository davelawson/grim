#!/usr/bin/env sh
set -eu
cd "$(dirname "$0")/sql"
sqlite3 -bail "$GRIM_DB" <create-database.sql
sqlite3 -bail "$GRIM_DB" <install-card-designs.sql
sqlite3 -bail "$GRIM_DB" <populate-test-data.sql
