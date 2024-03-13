.PHONY: all test clean

test:
	@echo "Ejecutando tests..."
	@go test -v ./... --coverprofile coverfile_out >> /dev/null
	@go tool cover -func coverfile_out

build:
	go build -mod=vendor -ldflags '-s -w' -o build/bin/process-loan cmd/main.go

mocks:
	go generate ./...

proto-generate:
	protoc --proto_path=internal/adapters/handlers/grpc/proto internal/adapters/handlers/grpc/proto/*.proto --go_out=plugins=grpc:internal/adapters/handlers/grpc/pb --grpc-gateway_out=:internal/adapters/handlers/grpc/pb

proto-clean:
	rm internal/adapters/handlers/grpc/pb/*.go

start:
	@cp .env ./configs/.env
	@go run ./cmd/main.go