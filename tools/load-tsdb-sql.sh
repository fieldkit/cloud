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

psql -h 127.0.0.1 -p 5433 -U postgres postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'fk';"
psql -h 127.0.0.1 -p 5433 -U postgres postgres -c "DROP DATABASE fk;" || true
psql -h 127.0.0.1 -p 5433 -U postgres postgres -c "CREATE DATABASE fk;"

if [ "${sql_file: -4}" == ".bz2" ]; then
	bunzip2 -c $sql_file | psql -h 127.0.0.1 -p 5433 -U postgres fk
elif [ "${sql_file: -3}" == ".xv" ]; then
	xv -c $sql_file | psql -h 127.0.0.1 -p 5433 -U postgres fk
else
	psql -h 127.0.0.1 -p 5433 -U postgres fk < $sql_file
fi

make migrate-up-tsdb
