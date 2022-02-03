NAME = serve

MAIN := ./cmd/$(NAME)/
BIN  := $(NAME)

all: clean build

build: fmt test
	go build -o $(BIN) $(MAIN)

test: fmt
	go test ./...

fmt:
	go fmt ./...

clean:
	go clean
	rm -f $(BIN)

.PHONY: all build test fmt clean
