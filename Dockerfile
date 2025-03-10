# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make api tools

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /bin/* .
EXPOSE 8080
CMD ["./agentAI-api"]
