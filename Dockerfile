FROM golang:1.8

RUN apt-get update && apt-get install -y zip

RUN curl -OL https://github.com/google/protobuf/releases/download/v3.2.0/protoc-3.2.0-linux-x86_64.zip &&\
	unzip protoc-3.2.0-linux-x86_64.zip -d protoc3 &&\
	cp protoc3/bin/* /usr/bin/

RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u google.golang.org/grpc

ENV PATH="/go/protoc3/bin:${PATH}"
