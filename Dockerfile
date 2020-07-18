FROM golang:1.14.6-alpine

WORKDIR /go/src/api

COPY . .

RUN apk add bash \
        curl \
        git \
    && curl -fLo /go/bin/air https://git.io/linux_air \
    && chmod +x /go/bin/air

RUN CGO_ENABLED=0 go build -o /go/bin/app

CMD air
