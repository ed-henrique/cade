FROM node:23.2.0-alpine AS builder

WORKDIR /app

COPY package.json .
COPY package-lock.json .
RUN npm install

COPY . .
RUN npm run build

FROM nginx:1.28-alpine as runner

COPY --from=builder /app/dist /usr/share/nginx/html
