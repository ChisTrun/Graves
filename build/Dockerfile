FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd
EXPOSE 8080

CMD ["./server", "-c", "./configs/config.yaml"]


