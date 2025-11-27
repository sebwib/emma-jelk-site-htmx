# Build stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache nodejs npm
RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app

COPY go.mod go.sum package*.json ./
RUN go mod download
RUN npm ci

# Copy all source code, inkl database.db
COPY . .

RUN templ generate
RUN npm run build
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
COPY --from=builder /app/database.db ./database.db
COPY --from=builder /app/database.db-shm ./database.db-shm
COPY --from=builder /app/database.db-wal ./database.db-wal

EXPOSE 8080
CMD ["./main"]
