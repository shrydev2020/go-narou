.DEFAULT_GOAL := build
.PHONY: build
build:
	go build -o narou

.PHONY: genProto
genProto:
	find ./usecase -name "*.proto" | xargs clang-format -i  -style=file
	find ./usecase -name "*.proto" | xargs \
	 protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative

.PHONY: fmtProto
fmtProto:
	find . -name "*.proto" | xargs clang-format -i  -style=file

.PHONY: help
	help:
		@echo "make file"