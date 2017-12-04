FROM alpine:3.5

LABEL maintainer "pinto.bikez@gmail.com"

ARG APP_NAME=stock-service

RUN apk add --no-cache ca-certificates

ADD ./build/$APP_NAME /app
ADD ./core.database.yml.example /core.database.yml
ADD ./queue.json /

# Environment Variables
ENV LISTEN "0.0.0.0:8080"
ENV DATABASE_FILE "core.database.yml"

CMD ["/app"]