FROM golang:1.17.1-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go mod tidy

RUN go build -o main

EXPOSE 8090

CMD ["./main"]