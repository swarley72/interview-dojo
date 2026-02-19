FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY proto/ ./proto/

COPY core-service/go.mod core-service/go.sum ./core-service/
WORKDIR /app/core-service
RUN go mod edit -replace github.com/swarley72/interview-dojo/proto=../proto && \
    go mod tidy && \
    go mod download

COPY core-service/ .
RUN go mod edit -replace github.com/swarley72/interview-dojo/proto=../proto && \
    go mod tidy && \
    CGO_ENABLED=0 go build -o /app/bin/core-service ./cmd/core-service

FROM alpine:3.21 AS runtime
COPY --from=builder /app/bin/core-service /usr/local/bin/core-service
CMD ["core-service"]
