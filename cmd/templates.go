package cmd

// --- JavaScript / TypeScript Ecosystem ---

// Vite Template (React, Vue, Svelte, etc.) - CSR
const viteDockerfile = `# Development stage
FROM node:22-alpine AS development
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 5173
CMD ["npm", "run", "dev", "--", "--host"]

# Build stage
FROM development AS builder
RUN npm run build

# Production stage
FROM nginx:alpine AS production
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
`

const viteCompose = `services:
  app:
    build:
      context: .
      target: development
    ports:
      - "5173:5173"
    volumes:
      - .:/app
      - /app/node_modules
    environment:
      - NODE_ENV=development
`

// Next.js Template - SSR
const nextDockerfile = `# Development stage
FROM node:22-alpine AS development
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 3000
# Ensure Host is 0.0.0.0 for Docker
CMD ["npm", "run", "dev", "--", "-H", "0.0.0.0"]

# Production stage
FROM node:22-alpine AS production
WORKDIR /app
COPY package*.json ./
RUN npm ci --omit=dev
COPY . .
RUN npm run build
CMD ["npm", "start"]
`

const nextCompose = `services:
  app:
    build:
      context: .
      target: development
    ports:
      - "3000:3000"
    volumes:
      - .:/app
      - /app/node_modules
    environment:
      - NODE_ENV=development
`

// Generic Node Template (Express, NestJS, etc.)
const nodeDockerfile = `# Development stage
FROM node:22-alpine AS development
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 3000
CMD ["npm", "run", "dev"]

# Production stage
FROM node:22-alpine AS production
WORKDIR /app
COPY package*.json ./
RUN npm ci --omit=dev
COPY . .
CMD ["npm", "start"]
`

const nodeCompose = `services:
  app:
    build:
      context: .
      target: development
    ports:
      - "3000:3000"
    volumes:
      - .:/app
      - /app/node_modules
    environment:
      - NODE_ENV=development
`

const nodeDockerignore = `node_modules
dist
.git
.github
.vscode
.next
`

// --- Generic Language Families ---

// Python Template (Flask, Django, FastAPI)
const pythonDockerfile = `# Development stage
FROM python:3.11-slim AS development
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
# Standard port for many Python frameworks
EXPOSE 8000
CMD ["python", "app.py"]

# Production stage
FROM python:3.11-slim AS production
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
# Hint: Use gunicorn for production if applicable
CMD ["python", "app.py"]
`

const pythonCompose = `services:
  app:
    build:
      context: .
      target: development
    ports:
      - "8000:8000"
    volumes:
      - .:/app
    environment:
      - FLASK_ENV=development
      - PYTHONUNBUFFERED=1
`

// Go Template (Generic)
const goDockerfile = `# Development stage
FROM golang:1.22-alpine AS development
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Install Air for hot reload
RUN go install github.com/air-verse/air@latest
CMD ["air"]

# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

# Production stage
FROM alpine:latest AS production
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]
`

const goCompose = `services:
  app:
    build:
      context: .
      target: development
    ports:
      - "8080:8080"
    volumes:
      - .:/app
      - /go/pkg/mod
`

// Java Template (Maven/Gradle generic-ish)
const javaDockerfile = `# Development stage
FROM maven:3.9-eclipse-temurin-21 AS development
WORKDIR /app
COPY pom.xml .
RUN mvn dependency:go-offline
COPY src ./src
CMD ["mvn", "spring-boot:run"]

# Production stage
FROM eclipse-temurin:21-jre-alpine AS production
WORKDIR /app
COPY --from=development /app/target/*.jar app.jar
CMD ["java", "-jar", "app.jar"]
`

const javaCompose = `services:
  app:
    build:
      context: .
      target: development
    ports:
      - "8080:8080"
    volumes:
      - .:/app
      - ~/.m2:/root/.m2
`

// PHP Template (Apache)
const phpDockerfile = `# Development stage
FROM php:8.2-apache AS development
WORKDIR /var/www/html
COPY composer.json composer.lock ./
RUN apt-get update && apt-get install -y unzip
RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer
RUN composer install
COPY . .
CMD ["apache2-foreground"]

# Production stage
FROM php:8.2-apache AS production
WORKDIR /var/www/html
COPY . .
RUN chown -R www-data:www-data /var/www/html
`

const phpCompose = `services:
  app:
    build:
      context: .
      target: development
    ports:
      - "80:80"
    volumes:
      - .:/var/www/html
`
