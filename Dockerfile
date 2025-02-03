# syntax=docker/dockerfile:1

# For building the apps/services.
FROM pkgxdev/pkgx:v2.2.0 AS base

ARG APP_SERVICE
ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
WORKDIR /go/src/app

# Install dependencies
RUN apt update -qq && \
    apt install --no-install-recommends -y curl && \
    rm -rf /var/lib/apt/lists/*

# Setup pkgx
COPY package.json ./
RUN dev

FROM base AS builder

# Install Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 go build -tags=release -o=dist/app ./backends/${APP_SERVICE}

# For running the apps/services on production.
FROM gcr.io/distroless/static-debian12

ARG APP_SERVICE
ENV PORT=3000
ENV APP_SERVICE=${APP_SERVICE}

COPY --from=builder /go/src/app/dist/app /

CMD ["/app", "start"]
