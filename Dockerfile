FROM node:20 AS node-builder
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

WORKDIR /app

COPY package.json ./
COPY pnpm-lock.yaml ./

RUN pnpm install

COPY ./web ./web

RUN npx golte dev

FROM golang:1.22.4 AS go-builder

WORKDIR /app

COPY --from=node-builder /app/build ./build

ENV GOCACHE=/root/.cache/go-build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN --mount=type=cache,target="/root/.cache/go-build" go build -ldflags "-s -w"

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=go-builder /app/pawste /root/pawste
COPY --from=go-builder /app/build /root/build
COPY ./web /root/web

EXPOSE 9454

CMD ["/root/pawste"]