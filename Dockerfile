FROM golang:1.13.15-alpine
LABEL maintainer="Peter Karacsonyi <peter.karacsonyi@msci.com>"

RUN apk add --no-cache git alpine-sdk

COPY src/*.go /go/
COPY db/* /go/

ENV REDIS_URL="redis://redis:6379"

RUN go get -d -v
RUN go build -o rep

EXPOSE  8000

ENTRYPOINT ["/go/rep"]
