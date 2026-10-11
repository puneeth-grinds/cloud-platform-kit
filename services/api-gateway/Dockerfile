# --- Stage 1: Build Stage ---

FROM golang:1.27.2-alpine3.24 AS build

WORKDIR /app

COPY go.mod .

RUN go mod download 

# --- Stage 2: Production Stage ---

