# syntax=docker/dockerfile:1

# Base stage for all builds
FROM jdxcode/mise:2025.8.20 AS base

ARG APP_SERVICE
ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
WORKDIR /go/src/app

# Setup mise
COPY package.json ./
COPY mise.toml ./
RUN mise settings experimental=true
RUN mise trust
# Handle both Docker secrets (CI/CD) and build args (local development)
# This token doesn't have any scope which is only used by mise to query/install
# public packages.
RUN --mount=type=secret,id=MISE_GITHUB_TOKEN,env=MISE_GITHUB_TOKEN mise install

# Install dependencies
RUN apt update -qq && \
    apt install --no-install-recommends -y apt-transport-https ca-certificates curl gnupg && \
    apt update -qq && \
    rm -rf /var/lib/apt/lists/*

# Backend builder stage
FROM base AS backend-builder

# Install Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 go build -tags=release -o=dist/app ./backends/${APP_SERVICE}

# Backend runtime stage
FROM gcr.io/distroless/static-debian12 AS backend

ARG APP_SERVICE
ENV PORT=3000
ENV APP_SERVICE=${APP_SERVICE}

COPY --from=backend-builder /go/src/app/dist/app /

CMD ["/app", "start"]

# SPA base stage
FROM jdxcode/mise:2025.8.20 AS spa-base

ARG APP_SERVICE
ARG VITE_API_BASE_URL
ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
WORKDIR /app

# Install dependencies
RUN apt update -qq && \
    apt install --no-install-recommends -y curl && \
    rm -rf /var/lib/apt/lists/*

# Setup mise
COPY package.json ./
COPY mise.toml ./
RUN mise settings experimental=true
RUN mise trust
# Handle both Docker secrets (CI/CD) and build args (local development)
# This token doesn't have any scope which is only used by mise to query/install
# public packages.
RUN --mount=type=secret,id=MISE_GITHUB_TOKEN,env=MISE_GITHUB_TOKEN mise install

# SPA builder stage
FROM spa-base AS spa-builder

# Pass the build argument to environment variable for RR7
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}

# Install all dependencies with cache mount for faster builds
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY ./apps/assets/package.json ./apps/assets/package.json
COPY ./apps/dashboard/package.json ./apps/dashboard/package.json
COPY ./packages/api/package.json ./packages/api/package.json
COPY ./packages/typescript-config/package.json ./packages/typescript-config/package.json
COPY ./packages/ui/package.json ./packages/ui/package.json
RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store \
    pnpm install --frozen-lockfile --prefer-offline

# Copy source and build
COPY . .

# Build the SPA application
RUN pnpm --filter=./apps/${APP_SERVICE} build

# SPA runtime stage
FROM nginx:alpine AS spa

ARG APP_SERVICE

# Copy the built assets
COPY --from=spa-builder /app/apps/${APP_SERVICE}/build/client /usr/share/nginx/html

# Configure nginx for React Router SPA
RUN echo 'error_log /dev/stderr warn;' > /etc/nginx/nginx.conf && \
    echo 'events { worker_connections 1024; }' >> /etc/nginx/nginx.conf && \
    echo 'http { \
    include /etc/nginx/mime.types; \
    default_type application/octet-stream; \
    error_log /dev/stderr warn; \
    access_log /dev/stdout combined; \
    include /etc/nginx/conf.d/*.conf; \
}' >> /etc/nginx/nginx.conf

RUN echo 'server { \
    listen 80; \
    server_name _; \
    root /usr/share/nginx/html; \
    \
    # Enable gzip compression \
    gzip on; \
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript; \
    \
    # RR7 assets files \
    location /assets/ { \
        alias /usr/share/nginx/html/assets/; \
        expires 365d; \
        access_log off; \
    } \
    \
    # Handle all other routes \
    location / { \
        try_files $uri $uri.html $uri/ /index.html; \
        add_header Cache-Control "no-cache, no-store, must-revalidate"; \
    } \
}' > /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
