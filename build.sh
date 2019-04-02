#!/bin/bash

USER_ID=`id -u $USER`

banner() {
    echo "+------------------------------------------+"
    printf "| %-40s |\n" "`date`"
    echo "|                                          |"
    printf "|`tput bold` %-40s `tput sgr0`|\n" "$@"
    echo "+------------------------------------------+"
}

function show_help {
    echo "build.sh [-h] [-p]"
    echo " -h this text"
    echo " -p push image"
}

# A POSIX variable
OPTIND=1
push=0

while getopts "h?p" opt; do
    case "$opt" in
        h|\?)
            show_help
            exit 0
            ;;
        p)  push=1
            ;;
    esac
done

shift $((OPTIND-1))

[ "$1" = "--" ] && shift

set -e
cd `dirname $0`

rm -rf build

banner "NODE-BASE"
docker build -t fk-server-node-base node-base

banner "SERVER"
docker build -t fk-server-build server

banner "ADMIN"
docker build -t fk-admin-build admin

banner "FRONTEND"
docker build -t fk-frontend-build frontend

banner "LANDING"
docker build -t fk-landing-build landing

banner "BUILDING"

mkdir build

mkdir build/tmp

mkdir build/api
docker rm -f fk-server-build > /dev/null 2>&1 || true
docker run --rm --name fk-server-build -v `pwd`/build:/build fk-server-build \
       sh -c "cp -r \$GOPATH/bin/server /build && cp -r api/public /build/api/ && chown -R $USER_ID /build/api /build/server"

mkdir build/admin
docker rm -f fk-admin-build > /dev/null 2>&1 || true
docker run --rm --name fk-admin-build -v `pwd`/build/admin:/build fk-admin-build \
       sh -c "cp -r /usr/app/build/* /build/ && chown -R $USER_ID /build"

mkdir build/frontend
docker rm -f fk-frontend-build > /dev/null 2>&1 || true
docker run --rm --name fk-frontend-build -v `pwd`/build/frontend:/build fk-frontend-build \
       sh -c "cp -r /usr/app/build/* /build/ && chown -R $USER_ID /build"

mkdir build/landing
docker rm -f fk-landing-build > /dev/null 2>&1 || true
docker run --rm --name fk-landing-build -v `pwd`/build/landing:/build fk-landing-build \
       sh -c "cp -r /usr/app/build/* /build/ && chown -R $USER_ID /build"

banner "Final Container"

for e in css csv html js json map svg txt; do
    find build -iname '*.$e' -exec gzip -k9 {} \;
done

cp /etc/ssl/certs/ca-certificates.crt build

echo '.DS_Store
Dockerfile' > build/.dockerignore

echo 'FROM scratch
ENV FIELDKIT_ADDR=:80
ENV FIELDKIT_ADMIN_ROOT=/admin
ENV FIELDKIT_FRONTEND_ROOT=/frontend
ENV FIELDKIT_LANDING_ROOT=/landing
COPY . /
ADD ca-certificates.crt /etc/ssl/certs/
EXPOSE 80
ENTRYPOINT ["/server"]' > build/Dockerfile

if [ "$1" = "production" ]; then
	DOCKER_TAG=latest
elif [ "$1" != "" ]; then
	DOCKER_TAG=$1
else
	DOCKER_TAG=`git rev-parse --abbrev-ref HEAD`
fi

docker build -t conservify/fk-cloud:$DOCKER_TAG build

if [ "$push" == "1" ]; then
    banner "Pushing!"
    docker push conservify/fk-cloud:$DOCKER_TAG
fi
