# Use the official Bun image
# see all versions at https://hub.docker.com/r/oven/bun/tags
FROM oven/bun:1 AS base
WORKDIR /app

# Install dependencies
FROM base AS deps
COPY package.json bun.lockb* ./
RUN bun install --frozen-lockfile

# Build the app
FROM base AS builder
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Set production environment
ENV NODE_ENV=production

# Build the Next.js application
RUN bun run build

# Use nginx to serve static files
FROM nginx:alpine AS runner
WORKDIR /usr/share/nginx/html

# Copy the static output from builder
COPY --from=builder /app/out ./

# Copy custom nginx config if needed
# COPY ./nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 3000

# Update nginx to listen on port 3000
RUN sed -i 's/listen\s*80;/listen 3000;/g' /etc/nginx/conf.d/default.conf

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
