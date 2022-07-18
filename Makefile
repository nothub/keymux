GOFLAGS = -race
MODNAME = $(shell go list -m)

$(MODNAME): clean lint
	go build $(GOFLAGS) -o $@

run: $(MODNAME)
	sudo ./$(MODNAME)

clean:
	go clean
	go mod tidy

lint:
	go vet

.PHONY: run clean lint
