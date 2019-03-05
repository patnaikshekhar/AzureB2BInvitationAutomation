FROM alpine
COPY ./static /static
COPY main /app
CMD ["./app"]