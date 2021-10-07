# [<ru>] Основной backend-образ с приложением и php-fpm.
FROM php:8.0-fpm-alpine as backend
WORKDIR /app

# [<ru>] Устанавливаем системные зависимости, необходимые php-модули и менеджер зависимостей – composer.
RUN set -x \
    && apk add --no-cache \
        icu-libs \
        bash \
    && apk add --no-cache --virtual .build-deps \
        icu-dev \
        autoconf \
        openssl \
        make \
        g++ \
    && docker-php-ext-install -j$(nproc) \
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
# [<ru>] Устанавливаем зависимости приложения (`dump-autoload` будем делать после копирования остальных файлов).
RUN composer install -n --no-dev --no-cache --no-ansi --no-autoloader --no-scripts --prefer-dist

# [<ru>] Копируем в образ все остальные файлы приложения.
COPY --chown=www-data:www-data . /app/
# [<ru>] Перегенерируем автозагрузчик классов.
RUN composer dump-autoload -n --optimize

# [<ru>] Копируем конфигурацию php-fpm.
COPY .werf/zzz-php-fpm-config.conf /usr/local/etc/php-fpm.d/zzz-php-fpm-config.conf

EXPOSE 9000

#############################################################################

# [<ru>] NGINX-образ с публичными файлами.
FROM nginx:stable-alpine as frontend
WORKDIR /www

# [<ru>] Копируем публичные файлы из backend образа.
COPY --from=backend /app/public /www

# [<ru>] Копируем конфигурацию NGINX.
COPY .werf/nginx.conf /etc/nginx/nginx.conf

EXPOSE 8080
