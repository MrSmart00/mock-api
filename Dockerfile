FROM golang:1.14.6-alpine

WORKDIR /go/src/api

COPY . .

RUN apk add bash \
        curl \
        git && \
    go get -u github.com/cosmtrek/air && \
    go build -o /go/bin/air github.com/cosmtrek/air

CMD ["air"]