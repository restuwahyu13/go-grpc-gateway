GW := gw
PB := pb

.PHONY: gorun
gorun:
	go run --race -v .

.PHONY: protoc-install
protoc-install:
ifdef type
ifeq (${type}, ${GW}):
	go install --race -v github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest \
  github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
else ifeq (${type}, ${PB}):
	go install --race -v google.golang.org/protobuf/cmd/protoc-gen-go@latest \
  google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
endif
endif

.PHONY: protoc-gen
protoc-gen:
	rm -rf ${PWD}/stubs/**/*
	\
	protoc --proto_path= ./protos/*.proto \
	--proto_path= ./google/api/*.proto \
	--plugin=${HOME}/go/bin/protoc-gen-go \
	--plugin=${HOME}/go/bin/protoc-gen-grpc-gateway \
	--go_out="${PWD}/stubs" \
	--grpc-gateway_out="${PWD}/stubs"
	\
	protoc --proto_path= ./protos/*.proto \
	--proto_path= ./google/api/*.proto \
	--plugin=${HOME}/go/bin/protoc-gen-go-grpc \
	--plugin=${HOME}/go/bin/protoc-gen-grpc-gateway \
	--go-grpc_out="${PWD}/stubs" \
	--grpc-gateway_out="${PWD}/stubs"