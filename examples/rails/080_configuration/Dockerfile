# [<ru>] Базовый образ. Используется и для сборки ассетов, и как основа для backend-образа.
FROM ruby:2.7 as base
WORKDIR /app

RUN apt-get update && apt-get install -y default-mysql-client

# [<ru>] Копируем в образ файлы, нужные для установки зависимостей приложения.
COPY Gemfile Gemfile.lock ./
# [<ru>] Уставливаем зависимости приложения.
RUN bundle install

#############################################################################

# [<ru>] Временный образ, в котором произойдет сборка ассетов.
FROM base as assets

# [<ru>] Устанавливаем NodeJS.
RUN curl -fsSL https://deb.nodesource.com/setup_14.x | bash -
RUN apt-get update && apt-get install -y nodejs

# [<ru>] Устанавливаем Yarn.
RUN npm install yarn --global

# [<ru>] Устанавливаем JS зависимости приложения.
COPY yarn.lock package.json ./
RUN yarn install

# [<ru>] Копируем в образ все остальные файлы приложения.
COPY . .

# [<ru>] Собираем ассеты.
RUN SECRET_KEY_BASE=NONE RAILS_ENV=production rails assets:precompile

#############################################################################

# [<ru>] Основной backend-образ с Rails и Puma.
FROM base as backend

# [<ru>] Копируем в образ все остальные файлы приложения.
COPY . .

# [<ru>] Добавляем webpack-манифест без ассетов.
COPY --from=assets /app/public/packs/manifest.json /app/public/packs/manifest.json

#############################################################################

# [<snippet frontend>]
# [<ru>] NGINX-образ с собранными ранее ассетами.
FROM nginx:stable-alpine as frontend
WORKDIR /www

# [<ru>] Копируем собранные ассеты из предыдушего сборочного образа.
COPY --from=assets /app/public /www
# [<endsnippet frontend>]
