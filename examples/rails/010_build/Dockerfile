FROM ruby:2.7.1
WORKDIR /app

# [<en>] Install system dependencies
# [<ru>] Уставливаем системные зависимости
RUN apt-get -o Acquire::Check-Valid-Until=false -o Acquire::Check-Date=false update -qq && apt-get install -y build-essential libpq-dev libxml2-dev libxslt1-dev curl default-mysql-client

# [<en>] Install application dependencies
# [<ru>] Устанавливаем зависимости приложения
COPY Gemfile /app/Gemfile
COPY Gemfile.lock /app/Gemfile.lock
RUN bundle install

# [<en>] Add all the files from our repo to the image, including app source code
# [<ru>] Добавляем в образ все остальные файлы из нашего репозитория, включая исходный код приложения
COPY . .
