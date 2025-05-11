.PHONY: build run clean proto

build:
	go build -o build/calculator cmd/api/main.go

run:
	go run cmd/api/main.go

clean:
	rm -f api/proto/*.pb.go
	rm -f build/calculator

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/calculator.proto