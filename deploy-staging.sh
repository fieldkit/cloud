#!/bin/sh
set -ex
cd `dirname $0`

deploy () {
	scp fieldkit-staging.service core@$1:~/
	ssh core@$1 '\
		sudo mv ~/fieldkit-staging.service /etc/systemd/system/fieldkit.service  &&\
		sudo systemctl daemon-reload &&\
		sudo systemctl stop fieldkit &&\
		sudo systemctl enable fieldkit &&\
		sudo systemctl start fieldkit'
}

if [ -z $1 ]; then
SERVERS=(fieldkit-server-staging-a.aws.fieldkit.team)
else
SERVERS=("${@:1}")
fi

for server in ${SERVERS[@]}; do
deploy $server
done
