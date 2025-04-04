# Stage 1 - builder
FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 👇 Critical: static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./main.go

# Stage 2 - minimal final image
FROM scratch

# Add CA certs (needed for HTTPS)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Add app binary
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
