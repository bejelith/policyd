FROM golang:alpine3.16 AS build

ENV GOPATH=""

COPY ./ .

RUN go build -o /cmd ./pkg/main

FROM alpine

COPY --from=build /cmd /cmd

ENTRYPOINT ["/cmd"]
