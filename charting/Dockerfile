FROM node:16.13.2-buster-slim

WORKDIR /usr/app
COPY ./package.json ./yarn.lock /usr/app/
RUN apt-get update \
    && apt-get install -qq build-essential libcairo2-dev libpango1.0-dev libjpeg-dev libgif-dev librsvg2-dev
RUN yarn install
COPY en.json /usr/app/
COPY server.js /usr/app/
COPY vega /usr/app/vega/
ADD static.env /etc/static.env

CMD [ "node", "server.js" ] 

