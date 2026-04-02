FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# /go/bin/migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/bin/api ./cmd/api

FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache tzdata

RUN addgroup -S appuser && adduser -S appuser -G appuser

# binarnik teper: v obraze vnutti: /app/api.exe
COPY --from=builder /app/bin/api ./api

COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

COPY ./migrations ./migrations

USER appuser

EXPOSE 8080