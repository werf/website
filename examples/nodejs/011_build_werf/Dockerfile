FROM node:14-stretch
WORKDIR /app
COPY . .

RUN apt update
RUN apt install -y tzdata locales
RUN npm ci

EXPOSE 3000
CMD ['node','/app/app.js']
