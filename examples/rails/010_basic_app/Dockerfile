FROM ruby:2.7
WORKDIR /app

# [<ru>] Копируем в образ файлы, нужные для установки зависимостей приложения
# [<en>] Copy the files needed to install the application dependencies into the image
COPY Gemfile Gemfile.lock ./
# [<ru>] Уставливаем зависимости приложения
# [<en>] Install the application dependencies
RUN bundle install

# [<ru>] Копируем в образ все остальные файлы приложения
# [<en>] Copy all other application files into the image
COPY . .
