tests:
	go test -v ./...

tidy:
	go mod tidy

run:
	CONFIG_PATH=config/config.yaml go run ./cmd/url-shortener

build:
	go build -o bin/url-shortener ./cmd/url-shortener