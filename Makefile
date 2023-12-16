run: build
	@./bin/notebase

build:
	@go build -o bin/notebase

test:
	@go test -v ./...