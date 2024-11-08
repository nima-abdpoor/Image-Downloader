# Use the official Golang image as a base
FROM golang:1.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main ./delivery/http

# Start a new stage from scratch
FROM ubuntu:22.04

# Install necessary libraries (glibc)
RUN apt-get update && apt-get install -y libc6 && apt-get install -y netcat && rm -rf /var/lib/apt/lists/*

# Install migrate CLI
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/

# Copy the binary, config file, and migration files
COPY --from=builder /app/main .
COPY config.yml .
COPY db/migration ./db/migration

# Set up environment variables (replace with actual values if needed)
ENV POSTGRES_USER=root
ENV POSTGRES_PASSWORD=secret
ENV POSTGRES_DB=GID
ENV POSTGRES_HOST=postgresql
ENV POSTGRES_PORT=5432

# Expose port 8000
EXPOSE 8000

# Script to wait for the database and then run migrations
CMD ./wait-for-it.sh "$POSTGRES_HOST:$POSTGRES_PORT" -- \
    psql -h $POSTGRES_HOST -U $POSTGRES_USER -c "CREATE DATABASE $POSTGRES_DB;" && \
    migrate -path ./db/migration -database "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable" up && \
    ./main