FROM golang:1.14.6-alpine as build

WORKDIR /go/src/api

COPY . .

RUN apk add git \
    && go build -o ./bin/mock

FROM alpine

WORKDIR /api

COPY --from=build /go/src/api/bin/mock .

RUN addgroup go \
  && adduser -D -G go go \
  && chown -R go:go /api/mock

CMD ["./mock"]