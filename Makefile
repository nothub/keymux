MOD_NAME = $(shell go list -m)
BIN_NAME = $(shell basename $(MOD_NAME))
VERSION  = $(shell git describe --tags --abbrev=0 --match v[0-9]* 2> /dev/null)
LDFLAGS  = -ldflags="-X '$(MOD_NAME)/buildinfo.Version=$(VERSION)'"
GOFLAGS = -race

out/$(BIN_NAME): $(shell ls go.mod go.sum *.go **/*.go)
	go build $(LDFLAGS) $(GOFLAGS) -o $@

README.txt: out/$(BIN_NAME)
	./out/$(BIN_NAME) --help > README.txt

.PHONY: clean
clean:
	go clean
	go mod tidy
	rm -rf out
