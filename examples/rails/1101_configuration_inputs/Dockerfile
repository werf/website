ARG RUBY_VERSION
FROM ruby:${RUBY_VERSION}

ARG CACHE_VERSION
RUN echo CACHE_VERSION="${CACHE_VERSION}"

WORKDIR /app

RUN apt-get -o Acquire::Check-Valid-Until=false -o Acquire::Check-Date=false update -qq && apt-get install -y build-essential libpq-dev libxml2-dev libxslt1-dev curl default-mysql-client

COPY Gemfile /app/Gemfile
COPY Gemfile.lock /app/Gemfile.lock

ARG GITHUB_TOKEN
RUN export BUNDLE_GITHUB__COM="${GITHUB_TOKEN}:x-oauth-basic" bundle install

COPY . .
