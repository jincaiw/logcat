# Build stage for frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm install
COPY web/ ./
RUN npm run build

# Build stage for Go binary
FROM golang:1.22-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /app/web/dist ./web/dist
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags="-s -w" -o logcat ./cmd/logcat

# Final minimal image
FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata sqlite-libs && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app
COPY --from=backend-builder /app/logcat .
COPY --from=backend-builder /app/configs/config.yaml ./configs/
COPY --from=backend-builder /app/web/dist ./web/dist

RUN mkdir -p /app/data

EXPOSE 8080 5140 5140/udp

VOLUME ["/app/data", "/app/configs"]

ENV LOGCAT_DATABASE_TYPE=sqlite
ENV LOGCAT_SQLITE_PATH=/app/data/logcat.db

HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
    CMD wget -qO- http://localhost:8080/healthz || exit 1

ENTRYPOINT ["./logcat"]
CMD ["--config", "configs/config.yaml"]