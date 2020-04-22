# TODO write better Makefile (add test, deps, linux build, etc.)
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOGET=$(GOCMD) get
GOCLEAN=$(GOCMD) clean
BINARY_NAME=bin/mybinary.exe

all: build
build:
		$(GOBUILD) -o $(BINARY_NAME) -v
run:	build
		./$(BINARY_NAME)
deps:
		$(GOGET) github.com/gookit/color
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
