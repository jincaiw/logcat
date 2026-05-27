.PHONY: build run clean dev test

build:
	cd web && npm install && npm run build
	go build -o logcat .

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