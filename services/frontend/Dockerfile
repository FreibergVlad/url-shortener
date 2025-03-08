# syntax=docker/dockerfile:1

FROM node:lts-alpine AS dev

WORKDIR /app

COPY services/frontend/package*.json ./
RUN npm ci

COPY services/frontend/ .

ENTRYPOINT ["npm", "run", "dev"]