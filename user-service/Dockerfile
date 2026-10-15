FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init --parseDependency -g cmd/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/user_app ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/user_app .

COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./user_app"]