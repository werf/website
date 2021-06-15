FROM ruby:2.7.1
WORKDIR /app

# Install system dependencies
RUN apt-get -o Acquire::Check-Valid-Until=false -o Acquire::Check-Date=false update -qq && apt-get install -y build-essential libpq-dev libxml2-dev libxslt1-dev curl libmagickwand-dev gsfonts

# Install app dependencies
COPY Gemfile /app/Gemfile
COPY Gemfile.lock /app/Gemfile.lock
RUN bundle install

# Add all the files from our repo into our image, including app source code
COPY . .
