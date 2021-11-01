# [<en>] The basic image. It is used both for assembling assets and as the basis for a backend image.
# [<ru>] Базовый образ. Используется и для сборки ассетов, и как основа для backend-образа.
FROM ruby:2.7 as base
WORKDIR /app

RUN apt-get update && apt-get install -y default-mysql-client

# [<en>] Copy the files needed to install the application dependencies into the image.
# [<ru>] Копируем в образ файлы, нужные для установки зависимостей приложения.
COPY Gemfile Gemfile.lock ./
# [<en>] Install the application dependencies.
# [<ru>] Устанавливаем зависимости приложения.
RUN bundle install

#############################################################################

# [<en>] The temporary image to build assets.
# [<ru>] Временный образ, в котором произойдет сборка ассетов.
FROM base as assets

# [<en>] Install Node.js.
# [<ru>] Устанавливаем NodeJS.
RUN curl -fsSL https://deb.nodesource.com/setup_14.x | bash -
RUN apt-get update && apt-get install -y nodejs

# [<en>] Install Yarn.
# [<ru>] Устанавливаем Yarn.
RUN npm install yarn --global

# [<en>] Install JS application dependencies.
# [<ru>] Устанавливаем JS зависимости приложения.
COPY yarn.lock package.json ./
RUN yarn install

# [<en>] Copy all other application files into the image.
# [<ru>] Копируем в образ все остальные файлы приложения.
COPY . .

# [<en>] Build assets.
# [<ru>] Собираем ассеты.
RUN SECRET_KEY_BASE=NONE RAILS_ENV=production rails assets:precompile

#############################################################################

# [<en>] The primary backend image with Rails and Puma.
# [<ru>] Основной backend-образ с Rails и Puma.
FROM base as backend

# [<en>] Copy all other application files into the image.
# [<ru>] Копируем в образ все остальные файлы приложения.
COPY . .

# [<en>] Add a webpack manifest (no assets).
# [<ru>] Добавляем webpack-манифест без ассетов.
COPY --from=assets /app/public/packs/manifest.json /app/public/packs/manifest.json

#############################################################################

# [<snippet frontend>]
# [<en>] Add an NIGINX image with the pre-built assets.
# [<ru>] NGINX-образ с собранными ранее ассетами.
FROM nginx:stable-alpine as frontend
WORKDIR /www

# [<en>] Copy the pre-built assets from the above image.
# [<ru>] Копируем собранные ассеты из предыдушего сборочного образа.
COPY --from=assets /app/public /www
# [<endsnippet frontend>]
