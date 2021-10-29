# [<ru>] Базовый образ. Используется и для сборки ассетов, и как основа для backend-образа.
FROM php:8.0-fpm-alpine as base
WORKDIR /app

# [<ru>] Устанавливаем системные зависимости, необходимые php-модули и composer.
RUN set -x \
    && apk add --no-cache \
        mysql-client \
        mariadb-connector-c \
        icu-libs \
        bash \
    && apk add --no-cache --virtual .build-deps \
        icu-dev \
        autoconf \
        openssl \
        make \
        g++ \
    && docker-php-ext-install -j$(nproc) \
        pdo_mysql \
        sockets \
        opcache \
        pcntl \
        intl \
        1>/dev/null \
    && docker-php-source delete \
    && apk del .build-deps \
    && php -r "readfile('http://getcomposer.org/installer');" | php -- --install-dir=/usr/bin/ --filename=composer

# [<ru>] Копируем в образ файлы, нужные для установки зависимостей приложения.
COPY ./composer.* /app/
# [<ru>] Устанавливаем зависимости приложения (`dump-autoload` запустим после копирования остальных файлов приложения).
RUN composer install -n --no-dev --no-cache --no-ansi --no-autoloader --no-scripts --prefer-dist

#############################################################################

# [<ru>] Основной backend-образ с приложением и php-fpm.
FROM base as backend

# [<ru>] Копируем в образ все остальные файлы приложения.
COPY --chown=www-data:www-data . /app/
# [<ru>] Перегенерируем автозагрузчик классов.
RUN composer dump-autoload -n --optimize

# [<ru>] Копируем конфигурацию php-fpm.
COPY .werf/zzz-php-fpm-config.conf /usr/local/etc/php-fpm.d/zzz-php-fpm-config.conf

EXPOSE 9000

#############################################################################

# [<snippet assets>]
# [<ru>] Временный образ, в котором произойдет сборка ассетов.
FROM base as assets

# [<ru>] Устанавливаем Node.js.
RUN apk add --no-cache nodejs npm

# [<ru>] Устанавливаем JS-зависимости приложения.
COPY package.json ./
RUN npm install

# [<ru>] Копируем в образ все остальные файлы приложения.
COPY . .

# [<ru>] Собираем ассеты.
RUN npm run prod
# [<endsnippet assets>]

#############################################################################

# [<snippet frontend>]
# [<ru>] NGINX-образ с публичными файлами.
FROM nginx:stable-alpine as frontend
WORKDIR /www

# [<ru>] Копируем публичные файлы из assets образа.
COPY --from=assets /app/public /www

EXPOSE 8080
# [<endsnippet frontend>]
