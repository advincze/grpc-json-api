
all: build


setup:
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get -u github.com/golang/protobuf/protoc-gen-go

gen:
	mkdir -p bin
	mkdir -p pkg/helloworld
	protoc \
	-I. \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--grpc-gateway_out=logtostderr=true:pkg/helloworld \
	--swagger_out=logtostderr=true:bin \
	--go_out=plugins=grpc:pkg/helloworld helloworld.proto

clean:
	rm -rf pkg/helloworld
	rm -rf bin

build: gen
	go build -o bin/grpc-api-srv ./cmd/grpc-api-srv
	go build -o bin/grpc-api-cli ./cmd/grpc-api-cli

run:
	./bin/grpc-api-srv

try-pb:
	./bin/grpc-api-cli alice

try-json:
	curl -XPOST ':8080/v1/hello' -d '{"name":"bob"}'

try: try-pb try-json