#!/bin/sh
set -ex
cd `dirname $0`

deploy () {
	scp services/fieldkit.service services/fieldkit-twitter.service core@$1:~/
	ssh core@$1 '\
		sudo mv ~/fieldkit.service /etc/systemd/system/fieldkit.service  &&\
		sudo mv ~/fieldkit-twitter.service /etc/systemd/system/fieldkit-twitter.service  &&\
		sudo systemctl daemon-reload &&\
		sudo systemctl stop fieldkit &&\
		sudo systemctl enable fieldkit &&\
		sudo systemctl start fieldkit &&\
		sudo systemctl stop fieldkit-twitter &&\
		sudo systemctl enable fieldkit-twitter &&\
		sudo systemctl start fieldkit-twitter'
}

if [ -z $1 ]; then
SERVERS=(fieldkit-server-a.aws.fieldkit.org)
else
SERVERS=("${@:1}")
fi

for server in ${SERVERS[@]}; do
deploy $server
done
