FROM ruby:2.7
WORKDIR /app

# [<en>] Copy the files needed to install the application dependencies into the image.
# [<ru>] Копируем в образ файлы, нужные для установки зависимостей приложения.
COPY Gemfile Gemfile.lock ./
# [<en>] Install the application dependencies.
# [<ru>] Уставливаем зависимости приложения.
RUN bundle install

# [<en>] Copy all other application files into the image.
# [<ru>] Копируем в образ все остальные файлы приложения.
COPY . .
