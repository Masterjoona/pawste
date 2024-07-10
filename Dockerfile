FROM node:20-slim AS base

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
COPY golte.config.ts ./
COPY svelte.config.js ./
COPY ./web /app/web

FROM base AS build
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod
RUN pnpm dlx golte

FROM golang:1.22.5 AS go-build

WORKDIR /app
ENV GOCACHE=/root/.cache/go-build
COPY --from=build /app/pkg/build /app/pkg/build
COPY ./pkg /app/pkg/
COPY main.go ./
COPY go.mod go.sum ./

RUN --mount=type=cache,target="/root/.cache/go-build" go build -o pawste "-ldflags=-s -w"
#RUN go build -o pawste "-ldflags=-s -w"

FROM alpine:latest

WORKDIR /app
COPY --from=go-build /app/pawste .

ENTRYPOINT ["./pawste"]
