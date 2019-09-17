BINARY=yabumi

VERSION := 0.1.6
DATE := `date +%FT%T%z`
TARGET := $(BINARY) $(BINARY).mac $(BINARY).exe

LDFLAGS=-ldflags "-w -s -X main.version=$(VERSION) -X main.date=$(DATE)"

default: $(BINARY)

build: $(BINARY).$(VERSION).linux_amd64 $(BINARY).$(VERSION).darwin_amd64 $(BINARY).$(VERSION).windows_amd64.exe $(BINARY).$(VERSION).linux_arm $(BINARY).$(VERSION).linux_arm64

$(BINARY): $(BINARY).go
	GO111MODULE=on go build $(LDFLAGS) -o $(BINARY)

$(BINARY).$(VERSION).linux_amd64: $(BINARY).go
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY).$(VERSION).linux_amd64

$(BINARY).$(VERSION).darwin_amd64: $(BINARY).go
	GO111MODULE=on GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY).$(VERSION).darwin_amd64
$(BINARY).$(VERSION).linux_arm: $(BINARY).go
	GO111MODULE=on GOOS=linux GOARCH=arm go build $(LDFLAGS) -o $(BINARY).$(VERSION).linux_arm

$(BINARY).$(VERSION).linux_arm64: $(BINARY).go
	GO111MODULE=on GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY).$(VERSION).linux_arm64

$(BINARY).$(VERSION).windows_amd64.exe: $(BINARY).go
	GO111MODULE=on GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY).$(VERSION).windows_amd64.exe

test:
	GO111MODULE=on go test

clean:
	rm -f $(BINARY) $(BINARY).*_amd64 $(BINARY).*_arm $(BINARY).*_arm64 *.exe

.PHONY: default clean test 
