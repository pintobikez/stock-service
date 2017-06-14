FROM alpine:3.5

LABEL maintainer "ricardo.pinto@dafiti.com.br"

ARG APP_NAME=stock-service

RUN apk add --no-cache ca-certificates

ADD ./build/$APP_NAME /app
ADD ./core.database.yml.example /core.database.yml
ADD ./queue.json /

# Environment Variables
ENV SS_LISTEN "0.0.0.0:8080"
ENV SS_DATABASE_FILE "core.database.yml"

CMD ["/app"]