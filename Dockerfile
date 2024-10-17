FROM golang:1.22.2-alpine as builder
WORKDIR /app
RUN apk add --no-cache gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o zg-data-guard cmd/zg-data-guard/main.go

FROM alpine:latest
WORKDIR /root/
COPY .env .env
COPY build.properties build.properties
COPY internal/database/migrations ./internal/database/migrations
COPY internal/database/sqls ./internal/database/sqls
COPY internal/webserver/templates ./internal/webserver/templates
COPY internal/database/connector/scripts/postgres ./internal/database/connector/scripts/postgres
COPY --from=builder /app/zg-data-guard .
EXPOSE 8081
CMD ["./zg-data-guard"]
