FROM golang:1.23.2-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o auth-service ./cmd/main.go
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/auth-service /usr/local/bin/auth-service
ENV DB_HOST=db 
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=root
ENV DB_NAME=authdb
CMD ["/usr/local/bin/auth-service"]
