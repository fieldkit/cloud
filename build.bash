#!/bin/bash
set -ex
cd `dirname $0`

rm -rf docker/fieldkit docker/admin docker/frontend

docker build -t fieldkit-build . 

rm -rf build
mkdir build
docker rm -f fieldkit-build || true
docker run --rm --name fieldkit-build -v `pwd`/build:/build fieldkit-build cp /go/bin/fieldkit /build/

cp -r build/fieldkit docker/

./admin/build.bash
cp -r admin/build docker/admin

./frontend/build.bash
cp -r frontend/build docker/frontend

cd docker
docker build -t ocrnyc/fieldkit .
docker push ocrnyc/fieldkit
