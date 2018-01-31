.PHONY: build install generate vet compile test doc fmt lint docker

default: test compile

compile:
	go build -v -o ../../bin/proto2gql ./main.go

fmt:
	go fmt ./*.go

test:
	go test ./*.go
