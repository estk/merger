.PHONY: pb all install test bench
all: install

install:
	go build

pb:
	protoc -I pb/ pb/merger.proto --go_out=plugins=grpc:pb

bench:
	go test -bench=. -benchmem -v -benchtime=5s

test:
	go test -v
