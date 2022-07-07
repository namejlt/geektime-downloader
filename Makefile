# Go parameters
GOCMD=GO111MODULE=on CGO_ENABLED=1 go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGENERATE=$(GOCMD) generate
GOMODTITY=$(GOCMD) mod tidy

all: build run

build:
	rm -rf dist/
	$(GOMODTITY)
	$(GOGENERATE)
	$(GOBUILD) -o dist/geektime-downloader main.go

buildf:
	rm -rf dist/
	mkdir -p dist
	$(GOBUILD) -o dist/geektime-downloader main.go

test:
	$(GOTEST) -v ./...

clean:
	rm -rf dist/

run:
	./dist/geektime-downloader -h

stop:
	pkill -f dist/geektime-downloader

fetch:
	ps aux | grep geektime-downloader