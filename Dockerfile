FROM golang:1.16rc1-alpine3.13 AS build
WORKDIR /mailpie
COPY ./ .
RUN go build github.com/da-coda/mailpie .

FROM alpine:3.13
EXPOSE 1025
EXPOSE 1143
EXPOSE 8000
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /mailpie/mailpie .
CMD ["/root/mailpie", "-config", "/root/mailpie.yml"]
