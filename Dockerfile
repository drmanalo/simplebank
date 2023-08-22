# Build stage
FROM golang:1.21-alpine3.18 AS builder
WORKDIR /db
COPY . .

# Run stage
FROM alpine:3.18
WORKDIR /db
COPY db/migration ./db/migration
COPY .env.template .env
COPY app.env.template ./app.env
