# Start with the official Golang image
FROM golang:1.17

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o app .

EXPOSE 8080

CMD ["./app"]
