FROM golang:1.21.4-alpine3.18 as builder
WORKDIR /csubscriptionbot
RUN apk update && apk upgrade --available && sync && apk add --no-cache --virtual .build-deps
COPY . .
RUN go build -ldflags="-w -s" .
FROM alpine:3.18.4
RUN apk update && apk upgrade --available && sync
COPY --from=builder /csubscriptionbot/csubscriptionbot /csubscriptionbot
ENTRYPOINT ["/csubscriptionbot"]
