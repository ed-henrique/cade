FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o main main.go

FROM alpine:3.21 AS runner

COPY --from=builder /app/main .
CMD ["/main"]
