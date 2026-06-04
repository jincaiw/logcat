.PHONY: build build-web build-linux run clean dev test lint mod docker-build docker-up docker-down mysql-run

build:
	cd web && npm install && npm run build
	go build -o logcat .

build-web:
	cd web && npm install && npm run build

build-linux:
	cd web && npm install && npm run build
	@if [ "$$(go env GOOS)" != "linux" ] && [ -z "$$CC" ]; then \
		echo "CGO Linux builds require Linux or a Linux C cross-compiler. Use the GitHub Actions release workflow for v0.1.0 builds."; \
		exit 1; \
	fi
	mkdir -p bin
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bin/logcat .

run:
	go run . --config configs/config.yaml

clean:
	rm -f logcat
	rm -rf web/dist

dev:
	go run . --config configs/config.yaml

test:
	go test ./... -v -cover

lint:
	golangci-lint run ./...

mod:
	go mod tidy
	go mod download

docker-build:
	docker build -t logcat:local .

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

mysql-run:
	LOGCAT_DATABASE_TYPE=mysql go run . --config configs/config.yaml
