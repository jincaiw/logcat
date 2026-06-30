# syntax=docker/dockerfile:1

FROM node:22-alpine AS frontend
WORKDIR /src/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend ./
RUN npm run build

FROM golang:1.25-alpine AS backend
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
COPY --from=frontend /src/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags="-s -w" -o /out/logcat .

FROM alpine:3.22
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata wget && \
    addgroup -S logcat && adduser -S logcat -G logcat && \
    mkdir -p /app/data /app/templates && chown -R logcat:logcat /app
COPY --from=backend /out/logcat /app/logcat
COPY templates /app/templates
ENV SYSLG_ALERT_DATA_DIR=/app/data \
    SYSLG_ALERT_TEMPLATES_DIR=/app/templates \
    LOGCAT_OPEN_BROWSER=0 \
    TZ=UTC
EXPOSE 8080 5140/udp 5140/tcp
VOLUME ["/app/data"]
USER logcat
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/healthz >/dev/null || exit 1
ENTRYPOINT ["/app/logcat"]
CMD ["8080"]
