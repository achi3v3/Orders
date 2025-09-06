# Stage 1: Builder
FROM golang:1.24.3-alpine3.22 AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

# Stage 2: Final image
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /build/main .
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
CMD ["./main"]