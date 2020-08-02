# Go parameters

GOCMD=go
GOBUILD=$(GOCMD) build
GOGET=$(GOCMD) get
GOCLEAN=$(GOCMD) clean

ifdef OS
	BINARY_NAME=bin/HideMetaData.exe
else
	BINARY_NAME=bin/HideMetaData
endif

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
