FROM node:18.17.0-buster-slim
ARG GIT_HASH=missing
ARG VERSION=missing

WORKDIR /app
COPY ./package*.json /app/
RUN apt-get update \
    && apt-get install -qq libcairo2-dev libpango1.0-dev libjpeg-dev libgif-dev librsvg2-dev
RUN npm install
COPY tsconfig.* en.json /app/
COPY *.ts /app/
COPY vega /app/vega/

RUN echo "export GIT_HASH=$GIT_HASH" > /etc/static.env
RUN echo "export FIELDKIT_VERSION=$VERSION" >> /etc/static.env
RUN cat /etc/static.env

ENV PATH="/app/node_modules/.bin:${PATH}"

CMD [ "ts-node", "server.ts" ] 

