FROM node:7.4.0
WORKDIR /usr/app
COPY ./package.json /usr/app/
RUN npm install
ADD webpack.* /usr/app/
ADD src/vendors /usr/app/src/vendors
RUN npm run-script dll
COPY . /usr/app/
# RUN npm run-script build
