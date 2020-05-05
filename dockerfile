FROM golang:alpine AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go build -o bamboo ./cmd


FROM alpine:latest
RUN apk add --no-cache bash

WORKDIR src

COPY --from=builder /app/bamboo .
COPY --from=builder /app/migrations ./migrations

ENTRYPOINT "./bamboo"



