version=0.0.1
date=$(shell date "+(%d/%m/%Y)")

build:
	@go build -o bin/carteado -ldflags '-s -w -X "main.appVersion=${version}" -X "main.appDate=${date}"' cmd/api/main.go
	@cp config.json.example bin/config.json

run: build
	@cd bin
	@carteado

test:
	go test -v ./...
