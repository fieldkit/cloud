#!/bin/bash
set -ex
cd `dirname $0`

docker build -t fieldkit-admin-build . 

rm -rf build
mkdir build
docker rm -f fieldkit-admin-build || true
docker run --rm --name fieldkit-admin-build -v `pwd`/build:/build fieldkit-admin-build bash -c 'cp -r /usr/app/dist/* /build/'
find build -name '*.css' -exec gzip -k9 {} \;
find build -name '*.html' -exec gzip -k9 {} \;
find build -name '*.js' -exec gzip -k9 {} \;
find build -name '*.json' -exec gzip -k9 {} \;
find build -name '*.txt' -exec gzip -k9 {} \;
