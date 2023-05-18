FROM node:18-alpine
ARG CACHE_VERSION
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --no-progress
COPY . .
EXPOSE 3000
CMD ["npm", "start"]
