BINARY = terraform-provider-centrifyvault
VERSION = 0.1.0

BASE_PATH = $(shell pwd)
BUILD_PATH = $(BASE_PATH)/build

SYSTEM_ARCH ?= amd64
SYSTEM_OS ?= linux
ifeq ($(OS),Windows_NT)
    SYSTEM_OS := windows
    ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
        SYSTEM_ARCH := amd64
    else
        ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
            SYSTEM_ARCH := amd64
        endif
        ifeq ($(PROCESSOR_ARCHITECTURE),x86)
            SYSTEM_ARCH := i386
        endif
    endif
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        SYSTEM_OS := linux
    endif
    ifeq ($(UNAME_S),Darwin)
        SYSTEM_OS := darwin
    endif
    UNAME_M := $(shell uname -m)
    ifeq ($(UNAME_M),x86_64)
        SYSTEM_ARCH := amd64
    endif
    ifneq ($(filter %86,$(UNAME_M)),)
        SYSTEM_ARCH := i386
    endif
    ifneq ($(filter arm%,$(UNAME_M)),)
        SYSTEM_ARCH := arm
    endif
endif


default: install

build:
	go build -o ./${BINARY}

install: build
	#mv ${BINARY} ~/.terraform.d/plugins
	mkdir -p ~/.terraform.d/plugins/$(SYSTEM_OS)_$(SYSTEM_ARCH)/; \
	mv ./${BINARY} ~/.terraform.d/plugins/$(SYSTEM_OS)_$(SYSTEM_ARCH)/

.PHONY: build