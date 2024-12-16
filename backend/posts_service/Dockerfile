FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o posts-service ./cmd/main.go
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/posts-service /usr/local/bin/posts-service
ENV DB_HOST=db 
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=root
ENV DB_NAME=postsdb
CMD ["/usr/local/bin/posts-service"]
