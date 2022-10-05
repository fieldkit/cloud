#!/bin/bash

set -xe

sql_file=$1

if [ -z $sql_file ]; then
	echo "sql file is required"
	exit 2
fi

if [ ! -f $sql_file ]; then
	echo "sql file is required"
	exit 2
fi

psql -h 127.0.0.1 -p 5432 -U fieldkit postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'fieldkit';"
psql -h 127.0.0.1 -p 5432 -U fieldkit postgres -c "DROP DATABASE fieldkit;" || true
psql -h 127.0.0.1 -p 5432 -U fieldkit postgres -c "CREATE DATABASE fieldkit;"

if [ "${sql_file: -4}" == ".bz2" ]; then
	bunzip2 -c $sql_file | psql -h 127.0.0.1 -p 5432 -U fieldkit fieldkit
elif [ "${sql_file: -3}" == ".xz" ]; then
	xz -dc $sql_file | psql -h 127.0.0.1 -p 5432 -U fieldkit fieldkit
elif [ "${sql_file: -4}" == ".sql" ]; then
	psql -h 127.0.0.1 -p 5432 -U fieldkit fieldkit < $sql_file
fi

make migrate-up
