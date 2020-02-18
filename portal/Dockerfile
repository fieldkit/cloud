FROM fk-server-node-base

WORKDIR /usr/app
COPY ./package.json ./yarn.lock /usr/app/

# --no-cache: download package index on-the-fly, no need to cleanup afterwards
# --virtual: bundle packages, remove whole bundle at once, when done
RUN apk --no-cache --virtual build-dependencies add \
    python \
    make \
    g++ \
    && yarn install \
    && apk del build-dependencies

COPY . /usr/app/

ENV PUBLIC_PATH /
RUN yarn build
