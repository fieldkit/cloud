FROM node:7.5.0
WORKDIR /usr/app
RUN npm install -g yarn
COPY ./package.json ./yarn.lock /usr/app/
RUN yarn install
COPY webpack.* /usr/app/
COPY src/vendors /usr/app/src/vendors
RUN npm run-script dll
COPY . /usr/app/
RUN npm run-script build
