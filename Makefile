VERSION_FILE=$(shell echo $$GOPATH)/src/github.com/michigan-com/newsfetch/VERSION
CURRENT_VERSION=$(shell cat $(VERSION_FILE))
MAJOR=$(word 1, $(subst ., , $(CURRENT_VERSION)))
MINOR=$(word 2, $(subst ., , $(CURRENT_VERSION)))
PATCH=$(word 3, $(subst ., , $(CURRENT_VERSION)))
VER ?= $(MAJOR).$(MINOR).$(shell echo $$(($(PATCH)+1)))
COMMITHASH=$(shell git rev-parse --short HEAD)

build:
	go build -ldflags "-X main.VERSION=$(CURRENT_VERSION) -X main.COMMITHASH=$(COMMITHASH)"

install:
	go install -ldflags "-X main.VERSION=$(CURRENT_VERSION) -X main.COMMITHASH=$(COMMITHASH)"

bump:
	echo $(VER) > $(VERSION_FILE)

test:
	go test ./...

get:
	go get -t ./...
