FROM golang:1.17-alpine

# Need this to build the app
ENV CGO_ENABLED 0

RUN mkdir / app
ADD . /app
WORKDIR /app
RUN go clean --modcache
RUN go build -o main .

EXPOSE 8080

CMD ["/app/main"]