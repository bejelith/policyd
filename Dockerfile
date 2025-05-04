FROM golang:1.24-alpine3.21 AS build
ENV GOARCH=""

COPY ./ .

RUN go build -o /cmd ./pkg/cmd

FROM alpine

COPY --from=build /cmd /cmd

ENTRYPOINT ["/cmd"]
