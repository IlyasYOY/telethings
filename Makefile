.PHONY: test vet build run install setup setup-remove clean

# Build directory
BIN_DIR := ./bin
BINARY_NAME := telethings
BINARY_PATH := $(BIN_DIR)/$(BINARY_NAME)

CMD_PATH := ./cmd/telethings

test:
	go test ./...

vet:
	go vet ./...

build: $(BIN_DIR)
	go build -o $(BINARY_PATH) $(CMD_PATH)

run: build
	./$(BINARY_PATH)

install:
	go install $(CMD_PATH)

setup:
	./setup.sh

setup-remove:
	./setup_remove.sh

clean:
	rm -rf $(BIN_DIR)

$(BIN_DIR):
	mkdir -p $(BIN_DIR)
