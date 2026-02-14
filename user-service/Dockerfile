FROM golang:1.25-alpine AS builder

WORKDIR /app


COPY proto/ ./proto/

COPY user-service/go.mod user-service/go.sum ./user-service/
WORKDIR /app/user-service
RUN go mod edit -replace github.com/swarley72/interview-dojo/proto=../proto && \
    go mod tidy && \
    go mod download

COPY user-service/ .
RUN go mod edit -replace github.com/swarley72/interview-dojo/proto=../proto && \
    go mod tidy && \
    CGO_ENABLED=0 go build -o /app/bin/user-service ./cmd/user-service

FROM alpine:3.21 AS runtime
COPY --from=builder /app/bin/user-service /usr/local/bin/user-service
CMD ["user-service"]
