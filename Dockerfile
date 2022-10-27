FROM golang:1.19-alpine
WORKDIR /go/src/go-fiber-api-docker
COPY . .
RUN go build -o bin/server main.go
CMD ["./bin/server"]