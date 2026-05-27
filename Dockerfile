FROM golang:1.25 AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o k8s-mcp .

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=builder /build/k8s-mcp /usr/local/bin/k8s-mcp
ENTRYPOINT ["k8s-mcp","-mode","http", "-apiKey","12345678"]
