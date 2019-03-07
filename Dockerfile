FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY ./static /static
COPY ./demo1_b2b /app
CMD ["./app"]