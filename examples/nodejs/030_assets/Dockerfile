FROM node:12-alpine as builder
WORKDIR /app

# [<en>] Copy the files needed to install the application dependencies into the image.
# [<ru>] Копируем в образ файлы, нужные для установки зависимостей приложения.
COPY package.json package-lock.json ./

# [<en>] Install the application dependencies.
# [<ru>] Устанавливаем зависимости приложения.
RUN npm ci

# [<en>] Copy all other application files into the image.
# [<ru>] Копируем в образ все остальные файлы приложения.
COPY . .

# [<en>] Build static asseets.
# [<ru>] Собираем статические файлы.
RUN npm run build

#############################################################################

FROM node:12-alpine as backend
WORKDIR /app

# [<en>] Copy the files needed to install the application dependencies into the image.
# [<ru>] Копируем в образ файлы, нужные для установки зависимостей приложения.
COPY package.json package-lock.json ./

# [<en>] Install the application dependencies.
# [<ru>] Устанавливаем зависимости приложения.
RUN npm ci --production

# [<en>] The app files.
# [<ru>] Файлы самого приложения.
COPY app.js ./
RUN mkdir dist
COPY --from=builder /app/dist/*.html ./dist/
COPY --from=builder /app/bin         ./bin/
COPY --from=builder /app/routes      ./routes/


#############################################################################

# Add an NGINX image with the pre-built assets.
FROM nginx:stable-alpine as frontend
WORKDIR /www

# Copy the pre-built assets from the above image.
COPY --from=builder /app/dist /www/static

# Copy the NGINX configuration.
COPY .werf/nginx.conf /etc/nginx/nginx.conf