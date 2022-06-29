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

make binaries

rm -rf build-containers

banner "NODE-BASE"
docker build -t fk-server-node-base node-base

banner "SERVER"
docker build -t fk-server-build server

banner "PORTAL"
docker build -t fk-portal-build portal

banner "BUILDING"

mkdir build-containers

mkdir build-containers/tmp

mkdir build-containers/api
docker rm -f fk-server-build > /dev/null 2>&1 || true
docker run --rm --name fk-server-build -v $WORKING_DIRECTORY/build-containers:/build fk-server-build \
       sh -c "cp -r /app/build/* /build && cp -r api/public /build/api/ && chown -R $USER_ID.$GROUP_ID /build"

mkdir build-containers/portal
docker rm -f fk-portal-build > /dev/null 2>&1 || true
docker run --rm --name fk-portal-build -v $WORKING_DIRECTORY/build-containers/portal:/build fk-portal-build \
       sh -c "cp -r /usr/app/build/* /build/ && chown -R $USER_ID.$GROUP_ID /build"

docker rm -f fk-server-build > /dev/null 2>&1 || true
docker rm -f fk-portal-build > /dev/null 2>&1 || true

banner "Final Container"

for e in css csv html js json map svg txt; do
    find build-containers -iname '*.$e' -exec gzip -k9 {} \;
done

cp /etc/ssl/certs/ca-certificates.crt build-containers

echo export GIT_HASH=`git rev-parse HEAD` > build-containers/static.env
echo export FIELDKIT_VERSION=$VERSION >> build-containers/static.env

echo '.DS_Store
Dockerfile' > build-containers/.dockerignore

echo 'FROM scratch
COPY . /
ADD ca-certificates.crt /etc/ssl/certs/
ADD static.env /etc/static.env
EXPOSE 80
ENTRYPOINT ["/server"]' > build-containers/Dockerfile

docker build -t conservify/fk-cloud:$DOCKER_TAG build-containers

cp build-containers/static.env charting

cd charting 
tar -czh --exclude='./node_modules' . | docker build -t conservify/fk-charting:$DOCKER_TAG -
cd ..
