.PHONY: run build fmt vet staticcheck test test-race test-race-cover clean check

run:
	go run ./cmd/server

build:
	go build -o bin/redis-clone ./cmd/server

fmt:
	go fmt ./...

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

test:
	go test ./...

test-race:
	go test -race ./...

test-race-cover:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

clean:
	rm -rf bin coverage.out

check: fmt vet staticcheck test-race
