#!/bin/bash
set -ex
cd `dirname $0`

deploy () {
	scp fieldkit-server.service core@$1:~/
	ssh core@$1 "set -ex \
		&& docker pull ocrnyc/fieldkit \
		&& sudo mv ~/fieldkit-server.service /etc/systemd/system/ \
		&& sudo systemctl daemon-reload \
		&& sudo systemctl stop fieldkit-server \
		&& sudo systemctl enable fieldkit-server \
		&& sudo systemctl start fieldkit-server"
}

SERVERS=(fieldkit-server-a.aws.fieldkit.org)
for server in ${SERVERS[@]}
do
	deploy $server
done
