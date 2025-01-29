FROM golang:1.23-alpine AS builder

RUN apk update && apk upgrade && apk add alpine-sdk

WORKDIR /app

COPY go.mod /app
COPY go.sum /app
RUN go mod download

ADD . /app
RUN make

FROM alpine:latest AS runner
COPY --from=builder /app/bin/staking-cli /usr/bin/staking-cli

ENTRYPOINT ["staking-cli", "--non-interactive"]
