#!/bin/sh
set -ex
cd `dirname $0`

rm -rf build

docker build -t fieldkit-server-build server
docker build -t fieldkit-admin-build admin
docker build -t fieldkit-frontend-build frontend
docker build -t fieldkit-landing-build landing

mkdir build

docker rm -f fieldkit-server-build || true
docker run --rm --name fieldkit-server-build -v `pwd`/build:/build fieldkit-server-build sh -c 'cp -r $GOPATH/bin/server /build/server'

mkdir build/admin
docker rm -f fieldkit-admin-build || true
docker run --rm --name fieldkit-admin-build -v `pwd`/build/admin:/build fieldkit-admin-build sh -c 'cp -r /usr/app/build/* /build/'

mkdir build/frontend
docker rm -f fieldkit-frontend-build || true
docker run --rm --name fieldkit-frontend-build -v `pwd`/build/frontend:/build fieldkit-frontend-build sh -c 'cp -r /usr/app/build/* /build/'

mkdir build/landing
docker rm -f fieldkit-landing-build || true
docker run --rm --name fieldkit-landing-build -v `pwd`/build/landing:/build fieldkit-landing-build sh -c 'cp -r /usr/app/build/* /build/'

find -E build -regex '.*\.(css|csv|html|js|json|map|svg|txt)' -exec gzip -k9 {} \;

echo 'FROM scratch
COPY . /
' > build/Dockerfile

docker build -t ocrnyc/fieldkit build
