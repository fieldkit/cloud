FROM node:7.7.4
WORKDIR /usr/app
COPY ./package.json ./yarn.lock /usr/app/
RUN yarn install
COPY . /usr/app/
RUN yarn dll
RUN yarn build
