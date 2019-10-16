#!/bin/bash

USER_ID=`id -u $USER`
WORKING_DIRECTORY=`pwd`

banner() {
    set +x
    echo "+------------------------------------------+"
    printf "| %-40s |\n" "`date`"
    echo "|                                          |"
    printf "| %-40s |\n" "$@"
    echo "+------------------------------------------+"
    set -x
}

function show_help {
    echo "build.sh [-h] [-p] [-w WORKING_DIRECTORY]"
    echo " -h this text"
    echo " -p push image"
    echo " -w absolute path from HOST to this working copy"
}

# A POSIX variable
OPTIND=1
push=0

while getopts "h?pw:" opt; do
    case "$opt" in
        h|\?)
            show_help
            exit 0
            ;;
        w)
            WORKING_DIRECTORY=$OPTARG
            ;;
        p)  push=1
            ;;
    esac
done

shift $((OPTIND-1))

[ "$1" = "--" ] && shift

set -xe
cd `dirname $0`

rm -rf build

banner "NODE-BASE"
docker build -t fk-server-node-base node-base

banner "SERVER"
docker build -t fk-server-build server

banner "PORTAL"
docker build -t fk-portal-build portal

banner "OCR-PORTAL"
docker build -t fk-ocr-portal-build ocr-portal

banner "LEGACY"
docker build -t fk-legacy-build legacy

banner "BUILDING"

mkdir build

mkdir build/tmp

mkdir build/api
docker rm -f fk-server-build > /dev/null 2>&1 || true
docker run --rm --name fk-server-build -v $WORKING_DIRECTORY/build:/build fk-server-build \
       sh -c "cp -r /app/build/* /build && cp -r api/public /build/api/ && chown -R $USER_ID /build/api /build/server"

mkdir build/portal
docker rm -f fk-portal-build > /dev/null 2>&1 || true
docker run --rm --name fk-portal-build -v $WORKING_DIRECTORY/build/portal:/build fk-portal-build \
       sh -c "cp -r /usr/app/build/* /build/ && chown -R $USER_ID /build"

mkdir build/ocr-portal
docker rm -f fk-ocr-portal-build > /dev/null 2>&1 || true
docker run --rm --name fk-ocr-portal-build -v $WORKING_DIRECTORY/build/ocr-portal:/build fk-ocr-portal-build \
       sh -c "cp -r /usr/app/build/* /build/ && chown -R $USER_ID /build"

mkdir build/legacy
docker rm -f fk-legacy-build > /dev/null 2>&1 || true
docker run --rm --name fk-legacy-build -v $WORKING_DIRECTORY/build/legacy:/build fk-legacy-build \
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
ENV FIELDKIT_PORTAL_ROOT=/portal
ENV FIELDKIT_OCR_PORTAL_ROOT=/ocr-portal
ENV FIELDKIT_LEGACY_ROOT=/legacy
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
