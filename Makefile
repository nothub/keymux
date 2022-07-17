EXENAME = $(shell go list -m)
GOFLAGS = -race

build: clean lint
	go build $(GOFLAGS) -o $(EXENAME)

run: build
	sudo ./$(EXENAME)

clean:
	go clean
	go mod tidy

lint:
	go vet

.PHONY: build run clean lint
