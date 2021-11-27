FROM golang:1.17.1-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main

EXPOSE 8081

EXPOSE 1234

EXPOSE 2112

CMD ["./main"]