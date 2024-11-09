build:
	@go build -o bin/blogAPI cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/blogAPI