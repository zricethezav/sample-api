FROM golang:1.10.0 AS build
MAINTAINER Zach Rice <zricezrice@gmail.com>

WORKDIR /go/src/github.com/zricethezav/gannet-market-api
COPY . .
RUN CGO_ENABLED=0 go build -o bin/gannet-market-api *.go

FROM alpine:3.7
COPY --from=build /go/src/github.com/zricethezav/gannet-market-api/bin/* /usr/bin/
ENTRYPOINT ["gannet-market-api"]
