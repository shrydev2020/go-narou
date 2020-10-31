.DEFAULT_GOAL := build
.PHONY: build
build:
	go build -o narou

.PHONY: help
	help:
		@echo "make file"