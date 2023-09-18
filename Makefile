#!/usr/bin/make -ef

NAME		= go-dht-bootstrap-peer
VERSION		= 0.0.2
GO_VERSION	= 1.20

BUILD_IMAGE	?=golang:$(GO_VERSION)-alpine
GO					?= go$(GO_VERSION)
PREFIX			?= /usr/local

IMAGE				= docker.io/bahner/$(NAME)
MODULE_NAME = github.com/bahner/$(NAME)

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

default: clean tidy $(NAME)

all: default image push

init: go.mod tidy

go.mod:
	$(GO) mod init $(MODULE_NAME)

tidy: go.mod
	$(GO) mod tidy

$(NAME): tidy
	$(GO) build -o $(NAME)

clean:
	rm -f $(NAME)

distclean: clean
	rm -f $(shell git ls-files --exclude-standard --others)

openwrt: GOOS=linux
openwrt: GOARCH=arm64
openwrt:
	$(GO) build -o $(NAME)
	

image:
	docker build \
	--build-arg BUILD_IMAGE=$(BUILD_IMAGE) \
	-t $(IMAGE) \
	.

push:
	docker push $(IMAGE)

release:
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)
	docker tag $(IMAGE) $(IMAGE):$(VERSION)
	docker push $(IMAGE):$(VERSION)

install:
	install -Dm755 $(NAME) $(DESTDIR)$(PREFIX)/bin/$(NAME)

.PHONY: default init tidy install clean distclean
