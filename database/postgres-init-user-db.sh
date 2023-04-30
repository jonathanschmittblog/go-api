#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER api PASSWORD 'api';
    CREATE DATABASE api;
    GRANT ALL PRIVILEGES ON DATABASE api TO api;
EOSQL
