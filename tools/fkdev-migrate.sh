#!/bin/bash

uname -a

source /etc/docker/compose/fieldkit.env

migrations/migrate.linux-amd64 --path migrations --database $FIELDKIT_POSTGRES_URL up
