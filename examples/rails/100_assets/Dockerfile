FROM ruby:2.7.1 as rails

WORKDIR /app

# Install system dependencies
RUN apt-get -o Acquire::Check-Valid-Until=false -o Acquire::Check-Date=false update -qq && apt-get install -y build-essential libpq-dev libxml2-dev libxslt1-dev curl default-mysql-client

# Install app dependencies
COPY Gemfile /app/Gemfile
COPY Gemfile.lock /app/Gemfile.lock
RUN bundle install

FROM rails as assets_builder

# Install NodeJS 14 and yarn
RUN curl -sL https://deb.nodesource.com/setup_14.x | bash -
RUN apt-get install -y nodejs
RUN npm install yarn --global

# Install assets dependencies
COPY yarn.lock package.json /app/
RUN yarn install

# Add all the files from our repo into our image, including app source code
COPY . .

# Prepare assets
RUN SECRET_KEY_BASE=NONE RAILS_ENV=production rails assets:precompile

FROM nginx:stable-alpine as frontend

# Copy all assets
COPY --from=assets_builder /app/public /www

# Add nginx configuration
COPY .werf/nginx.conf /etc/nginx/nginx.conf

FROM rails as backend

# Add all the files from our repo into our image, including app source code
COPY . .

# Add webpack manifest (precompiled assets not needed)
COPY --from=assets_builder /app/public/packs/manifest.json /app/public/packs/manifest.json
