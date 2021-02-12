FROM ruby:2.7.1
WORKDIR /app
COPY . .

RUN apt-get update -qq && apt-get install -y \
 build-essential libpq-dev libxml2-dev libxslt1-dev curl
RUN bundle install

CMD ["bundle", "exec", "puma"]
