package cmd

const reactDockerfile = `# Development stage
FROM node:20-alpine AS development
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 5000
CMD ["npm", "run", "dev", "--", "--host", "--port", "5000"]

# Build stage
FROM development AS builder
RUN npm run build

# Production stage
FROM node:20-alpine AS production
RUN npm install -g serve
WORKDIR /app
COPY --from=builder /app/dist ./dist
EXPOSE 5000
CMD ["serve", "-s", "dist", "-l", "5000"]
`

const reactCompose = `services:
  app:
    build:
      context: .
      target: development
    ports:
      - "5000:5000"
    volumes:
      - .:/app
      - /app/node_modules
    environment:
      - NODE_ENV=development
`

const reactDockerignore = `node_modules
dist
.git
.github
.vscode
`
