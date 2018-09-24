BINARY=yabumi

VERSION := 0.1.3
DATE := `date +%FT%T%z`
TARGET := $(BINARY) $(BINARY).mac $(BINARY).exe

LDFLAGS=-ldflags "-w -s -X main.version=$(VERSION) -X main.date=$(DATE)"

default: $(BINARY)

build: $(BINARY).$(VERSION).linux_amd64 $(BINARY).$(VERSION).darwin_amd64 $(BINARY).$(VERSION).windows_amd64.exe

$(BINARY): $(BINARY).go
	go build $(LDFLAGS) -o $(BINARY)

$(BINARY).$(VERSION).linux_amd64: $(BINARY).go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY).$(VERSION).linux_amd64

$(BINARY).$(VERSION).darwin_amd64: $(BINARY).go
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY).$(VERSION).darwin_amd64

$(BINARY).$(VERSION).windows_amd64.exe: $(BINARY).go
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY).$(VERSION).windows_amd64.exe

test:
	go test

clean:
	rm -f $(BINARY) $(BINARY).*_amd64 *.exe

.PHONY: default clean test 
