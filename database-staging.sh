#!/bin/bash
exec ssh -t core@fieldkit-server-staging-a.aws.fieldkit.team docker run -it --rm postgres psql postgres://fieldkit@fieldkit-staging.chgdsru0ljqu.us-east-1.rds.amazonaws.com:5432/fieldkit?sslmode=require $@
