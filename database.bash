#!/bin/bash
exec ssh -t core@fieldkit-server-a.aws.fieldkit.org docker run -it --rm postgres psql postgres://fieldkit@fieldkit.chgdsru0ljqu.us-east-1.rds.amazonaws.com:5432/fieldkit?sslmode=require $@
