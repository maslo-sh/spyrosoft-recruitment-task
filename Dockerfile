FROM golang:1.18-bullseye AS build

WORKDIR /app

RUN apt-get update \
  && apt-get install gcc

COPY go.mod ./
RUN go mod download

COPY base ./base
COPY marshal ./marshal
COPY logger ./logger
COPY *.go ./

RUN go build -ldflags '-linkmode external -w -extldflags "-static"' -o /nbp-api-query-worker

FROM alpine:latest

WORKDIR /

COPY --from=build /nbp-api-query-worker /nbp-api-query-worker

ENTRYPOINT ["/nbp-api-query-worker"]