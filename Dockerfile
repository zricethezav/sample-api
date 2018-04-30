FROM golang:1.10.0 AS build
MAINTAINER Zach Rice <zricezrice@gmail.com>

WORKDIR /go/src/github.com/zricethezav/sample-api
COPY . .
RUN CGO_ENABLED=0 go build -o bin/sample-api *.go

FROM alpine:3.7
COPY --from=build /go/src/github.com/zricethezav/sample-api/bin/* /usr/bin/
ENTRYPOINT ["gannet-market-api"]
