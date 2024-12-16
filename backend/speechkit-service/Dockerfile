FROM golang:1.23.2-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o speechkit-service ./cmd/main.go
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/speechkit-service /usr/local/bin/speechkit-service
CMD ["/usr/local/bin/speechkit-service"]
