.PHONY: all build gen-pb gen-pb-gw
all: build
dep:
	dep ensure
build: dep gen-pb gen-pb-gw
	go build -o server/server github.com/hakobe/grpc-gateway-example/server
	go build -o gateway/gateway github.com/hakobe/grpc-gateway-example/gateway
gen-pb:
	protoc -I/usr/local/include -I. \
		-I$(GOPATH)/src \
		-I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:articles \
		./articles.proto
gen-pb-gw: 
	protoc -I/usr/local/include -I. \
		-I$(GOPATH)/src \
		-I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:articles \
		./articles.proto
run:
	./run.sh
