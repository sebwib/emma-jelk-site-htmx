# Build stage
FROM golang:1.24-alpine AS builder

# Install Node.js and npm
RUN apk add --no-cache nodejs npm

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum package*.json ./

# Install dependencies
RUN go mod download
RUN npm ci

# Copy source code
COPY . .

# Generate templ files
RUN templ generate

# Build Tailwind CSS
RUN npm run build

# Build Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary and static files from builder
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static

# Create empty db.db file (will be replaced by volume in production)
RUN touch db.db

EXPOSE 8080

CMD ["./main"]
