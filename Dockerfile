FROM golang:alpine AS builder

# Set the working directory
WORKDIR /app


# Copy the source code
COPY go.mod ./
COPY go.sum ./

# Download the packages
RUN go mod download

# Copy all the go files
COPY *.go ./

# Build the app
RUN go build -ldflags="-s -w" -o main

# RUN stage
FROM alpine:latest
WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/main . 
COPY postgres/migrations ./postgres/migrations
COPY app.env .

# Expose the app's port
EXPOSE 8000


CMD ["./main"]