.PHONY: build test

build:
	./build.sh

test:
	protoc --plugin=protoc-gen-gsk=../aa/protoc-gen-gsk --proto_path=. --go_out=../aa --gsk_out=../aa --gsk_opt=aa:bb:cc greeter.proto
