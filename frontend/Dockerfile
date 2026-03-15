FROM node:20.12.2
RUN apt-get update && \
    apt-get -y install default-jre
RUN npm i -g @angular/cli @openapitools/openapi-generator-cli

COPY . /frontend

WORKDIR /frontend

RUN npm ci

EXPOSE 4200
