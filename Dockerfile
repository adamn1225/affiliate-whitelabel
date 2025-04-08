# Stage 1 - builder
FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./main.go

# Stage 2 - development image with Go runtime
FROM golang:1.23

WORKDIR /app

# Copy the built binary and source code
COPY --from=builder /app/main .
COPY . .

# Install dependencies for running Go scripts
RUN go mod tidy

EXPOSE 8080

# Default command
CMD ["./main"]