FROM node:14-stretch
WORKDIR /app
COPY . .

RUN apt update
RUN apt install -y tzdata locales
RUN npm ci

CMD ['node','/app/app.js']
