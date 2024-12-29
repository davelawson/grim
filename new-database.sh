cd sql
sqlite3 $GRIM_DB <create-database.sql
sqlite3 $GRIM_DB <populate-test-data.sql
cd ..
