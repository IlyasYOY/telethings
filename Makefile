.PHONY: test coverage vet build run install setup cleanup clean

# Build directory
BIN_DIR := ./bin
BINARY_NAME := telethings
BINARY_PATH := $(BIN_DIR)/$(BINARY_NAME)

CMD_PATH := ./cmd/telethings

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out -coverpkg=./... ./...
	go tool cover -func=coverage.out \
		| sort \
		| sed "s|github.com/IlyasYOY/telethings|$(shell pwd)|g"

vet:
	go vet ./...

build: $(BIN_DIR)
	go build -o $(BINARY_PATH) $(CMD_PATH)

run: build
	./$(BINARY_PATH)

install:
	go install $(CMD_PATH)

generate:
	go generate ./...

setup:
	./setup.sh

cleanup:
	./cleanup.sh

clean:
	rm -rf $(BIN_DIR)

$(BIN_DIR):
	mkdir -p $(BIN_DIR)
