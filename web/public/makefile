.DEFAULT_GOAL := build
.PHONY: build
build:
	npx webpack

.PHONY: help
	help:
		@echo "make file"

.PHONY: genProto
genProto:
	find ./protos -name "*.proto" | xargs clang-format -i  -style=file
	find ./protos -name "*.proto" | xargs \
	 protoc \
		--ts_out=import_style=commonjs,binary:./src/protos \
		--plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts
