FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY proto/ ./proto/

COPY api-gateway/go.mod api-gateway/go.sum ./api-gateway/
WORKDIR /app/api-gateway
RUN go mod edit -replace github.com/swarley72/interview-dojo/proto=../proto && \
    go mod tidy && \
    go mod download

COPY api-gateway/ .
RUN go mod edit -replace github.com/swarley72/interview-dojo/proto=../proto && \
    go mod tidy && \
    CGO_ENABLED=0 go build -o /app/bin/api-gateway ./cmd/api-gateway

FROM alpine:3.21 AS runtime
COPY --from=builder /app/bin/api-gateway /usr/local/bin/api-gateway
CMD ["api-gateway"]
