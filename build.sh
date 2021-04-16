#!/bin/bash

if [ -z "$USER_ID" ]; then
	USER_ID=`id -u $USER`
	GROUP_ID=`id -gr`
fi

if [ -z "$WORKING_DIRECTORY" ]; then
	WORKING_DIRECTORY=`pwd`
fi

if [ -z "$DOCKER_TAG" ]; then
	DOCKER_TAG="main"
fi

banner() {
    set +x
    echo "+------------------------------------------+"
    printf "| %-40s |\n" "`date`"
    echo "|                                          |"
    printf "| %-40s |\n" "$@"
    echo "+------------------------------------------+"
    set -x
}

set -xe
cd `dirname $0`

rm -rf build

banner "NODE-BASE"
docker build -t fk-server-node-base node-base

banner "SERVER"
docker build -t fk-server-build server

banner "PORTAL"
docker build -t fk-portal-build portal

banner "BUILDING"

mkdir build

mkdir build/tmp

mkdir build/api
docker rm -f fk-server-build > /dev/null 2>&1 || true
docker run --rm --name fk-server-build -v $WORKING_DIRECTORY/build:/build fk-server-build \
       sh -c "cp -r /app/build/* /build && cp -r api/public /build/api/ && chown -R $USER_ID.$GROUP_ID /build"

mkdir build/portal
docker rm -f fk-portal-build > /dev/null 2>&1 || true
docker run --rm --name fk-portal-build -v $WORKING_DIRECTORY/build/portal:/build fk-portal-build \
       sh -c "cp -r /usr/app/build/* /build/ && chown -R $USER_ID.$GROUP_ID /build"

docker rm -f fk-server-build > /dev/null 2>&1 || true
docker rm -f fk-portal-build > /dev/null 2>&1 || true

banner "Final Container"

for e in css csv html js json map svg txt; do
    find build -iname '*.$e' -exec gzip -k9 {} \;
done

cp /etc/ssl/certs/ca-certificates.crt build

echo export GIT_HASH=`git rev-parse HEAD` > build/static.env
echo export FIELDKIT_VERSION=$VERSION >> build/static.env

echo '.DS_Store
Dockerfile' > build/.dockerignore

echo 'FROM scratch
COPY . /
ADD ca-certificates.crt /etc/ssl/certs/
ADD static.env /etc/static.env
EXPOSE 80
ENTRYPOINT ["/server"]' > build/Dockerfile

docker build -t conservify/fk-cloud:$DOCKER_TAG build
