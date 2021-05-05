.PHONY: build build-docker test fmt install uninstall

default: build

BINARY=go-clean-docker-registry

build:
	@echo ">> building amd64 binary"
	mkdir -p build
	go build -v -o build/ -v ./...

build-docker:
	@echo ">> building amd64 binary using docker"
	docker run --rm -v $(PWD):/usr/src/$(BINARY) \
      -w /usr/src/$(BINARY) \
      -e GOOS=linux \
      -e GOARCH=amd64 \
      golang:1.16 go get -d -v ./... && go install -v ./... && go build -v -o build/ -v ./...

test:
	go test -cover ./...

fmt:
	go fmt ./...

install:
	chmod 755 build/$(BINARY)
	sudo cp build/$(BINARY) /usr/local/bin/
	#todo config file
	#sudo mkdir /etc/$(BINARY)
	#sudo cp etc/config.yml /etc/$(BINARY)/config.yml

uninstall:
	sudo rm /usr/local/bin/$(BINARY)
