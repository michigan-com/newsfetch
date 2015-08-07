VERSION_FILE=$(shell echo $$GOPATH)/src/github.com/michigan-com/newsfetch/VERSION

build:
	go build -ldflags "-X main.VERSION $(shell cat $(VERSION_FILE))"
