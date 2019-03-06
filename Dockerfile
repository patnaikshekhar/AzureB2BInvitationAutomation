FROM alpine
COPY ./static /static
COPY ./demo1_b2b /app
CMD ["./app"]