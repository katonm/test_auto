FROM golang:1.14 as builder

COPY . /app
WORKDIR /app/cmd/

RUN CGO_ENABLED=0 GOOS=linux go build -o app

ENTRYPOINT ["./app"]
