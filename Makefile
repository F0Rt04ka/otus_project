
build-client:
	go build -o client.out ./cmd/client/main.go

build-daemon:
	go build -o daemon.out ./cmd/daemon/main.go

generate:
	protoc -I proto proto/system_monitor.proto --go_out=./proto/gen/go/ --go_opt=paths=source_relative --go-grpc_out=./proto/gen/go/ --go-grpc_opt=paths=source_relative

tests: 
	go test ./test
	go test --race -v ./...