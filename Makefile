version=0.0.1
date=$(shell date "+(%d/%m/%Y)")

.PHONY: frontend build

frontend:
	cd frontend && npm run build && cp -r dist ../cmd/server

build:
	go build -o bin/carteado -ldflags "-s -w -X 'main.appVersion=${version}' -X 'main.appDate=${date}'" cmd/server/main.go

all: frontend build-win
	@echo "all is done now!"

run:
	cd bin
	carteado
	cd ..

test:
	go test -v ./...
