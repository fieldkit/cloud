#!/bin/bash
set -ex
cd `dirname $0`

docker build -t fieldkit-frontend-build . 

rm -rf build
mkdir build
docker rm -f fieldkit-frontend-build || true
docker run --rm --name fieldkit-frontend-build -v `pwd`/build:/build fieldkit-frontend-build bash -c 'cp -r /usr/app/dist/* /build/'
find build -name '*.css' -exec gzip -k9 {} \;
find build -name '*.html' -exec gzip -k9 {} \;
find build -name '*.js' -exec gzip -k9 {} \;
find build -name '*.json' -exec gzip -k9 {} \;
find build -name '*.txt' -exec gzip -k9 {} \;
