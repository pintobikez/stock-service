FROM alpine:3.5

LABEL maintainer "fabio.ribeiro@dafiti.com.br"

ARG APP_NAME=stock-service

RUN apk add --no-cache ca-certificates

ADD ./build/$APP_NAME /app
ADD ./core.database.yml /

CMD ["/app"]
