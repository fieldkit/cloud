FROM node:7.9.0-alpine
WORKDIR /usr/app
COPY ./package.json ./yarn.lock /usr/app/
RUN yarn install
COPY . /usr/app/
RUN yarn build
