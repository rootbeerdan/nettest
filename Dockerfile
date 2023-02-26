FROM registry.ipv6.docker.com/library/golang:alpine

WORKDIR /app

COPY . .

RUN go build -o nettest

EXPOSE 8081

ENTRYPOINT ["./nettest"]
