FROM golang:1.24.3-alpine3.22

WORKDIR /producer-service

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

CMD ["./main"]