#!/bin/bash

USER_ID=`id -u $USER`

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

set -ex
cd `dirname $0`

rm -rf build

docker build -t fieldkit-server-build server
docker build -t fieldkit-admin-build admin
docker build -t fieldkit-frontend-build frontend
docker build -t fieldkit-landing-build landing

mkdir build

mkdir build/tmp

mkdir build/api
docker rm -f fieldkit-server-build > /dev/null 2>&1 || true
docker run --rm --name fieldkit-server-build -v `pwd`/build:/build fieldkit-server-build \
       sh -c "cp -r \$GOPATH/bin/server /build && cp -r api/public /build/api/ && chown -R $USER_ID /build/api /build/server"

mkdir build/admin
docker rm -f fieldkit-admin-build > /dev/null 2>&1 || true
docker run --rm --name fieldkit-admin-build -v `pwd`/build/admin:/build fieldkit-admin-build \
       sh -c "cp -r /usr/app/build/* /build/ && chown -R $USER_ID /build"

mkdir build/frontend
docker rm -f fieldkit-frontend-build > /dev/null 2>&1 || true
docker run --rm --name fieldkit-frontend-build -v `pwd`/build/frontend:/build fieldkit-frontend-build \
       sh -c "cp -r /usr/app/build/* /build/ && chown -R $USER_ID /build"

mkdir build/landing
docker rm -f fieldkit-landing-build > /dev/null 2>&1 || true
docker run --rm --name fieldkit-landing-build -v `pwd`/build/landing:/build fieldkit-landing-build \
       sh -c "cp -r /usr/app/build/* /build/ && chown -R $USER_ID /build"

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
    docker push conservify/fk-cloud:$DOCKER_TAG
fi
