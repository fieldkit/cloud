#!/bin/sh

UID=`id -u $USER`

set -ex
cd `dirname $0`

rm -rf build

docker build -t fieldkit-landing-build landing

mkdir build

mkdir build/landing
docker rm -f fieldkit-landing-build > /dev/null 2>&1 || true
docker run --rm --name fieldkit-landing-build -v `pwd`/build/landing:/build fieldkit-landing-build \
       sh -c "cp -r /usr/app/build/* /build/ && chown -R $UID /build"

find build -regextype posix-extended -regex '.*\.(css|csv|html|js|json|map|svg|txt)' -exec gzip -k9 {} \;

echo 'Dockerfile' > build/.dockerignore

echo 'FROM nginx
COPY landing/ /usr/share/nginx/html/
RUN ls -alh /usr/share/nginx/html' > build/Dockerfile

if [ "$1" = "production" ]; then
	DOCKER_TAG=latest
elif [ "$1" != "" ]; then
	DOCKER_TAG=$1
else
	DOCKER_TAG=`git rev-parse --abbrev-ref HEAD`
fi

docker build -t conservify/fk-landing:$DOCKER_TAG build

docker push conservify/fk-landing:$DOCKER_TAG
