all: run

build:
	@go build -o ./bin/todoapi ./cmd/api/main.go

run: build
	@./bin/todoapi --watch