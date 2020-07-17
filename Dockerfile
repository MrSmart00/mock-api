FROM golang:1.14.6-alpine

WORKDIR /go/src/app

COPY . .

RUN apk add curl \
        git \
    && curl -fLo /go/bin/air https://git.io/linux_air \
    && chmod +x /go/bin/air

CMD air
