# Stage 1: Base image with pnpm installed
FROM node:20-alpine AS base

RUN npm i -g pnpm

# Stage 2: Install dependencies
FROM base AS dependencies

WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN pnpm install
COPY golte.config.ts ./
COPY ./web /app/web
RUN pnpm dlx golte

FROM golang:1.22.2 AS go-build

WORKDIR /app
#ENV GOCACHE=/root/.cache/go-build
COPY --from=dependencies /app/pkg/build /app/pkg/build
COPY ./pkg /app/pkg/
COPY main.go ./
COPY go.mod go.sum ./

#RUN --mount=type=cache,target="/root/.cache/go-build" go build -o pawste "-ldflags=-s -w"

FROM alpine:latest

WORKDIR /app
COPY --from=go-build /app/pawste .

ENTRYPOINT ["./pawste"]
