GRPC_DEPS:=\
	google.golang.org/protobuf/cmd/protoc-gen-go\
	google.golang.org/grpc/cmd/protoc-gen-go-grpc
run:
	@go run ./cmd/client/main.go

docker-generate:
	@docker run --rm \
		-v$(CURDIR):/generate \
		dafaque/go-grpc-generator:0.0.4-alpine3.16
	@go mod tidy

generate:
	@protoc --proto_path=api/protobuf\
			--go_out=internal/pb\
			--go_opt=paths=import\
			--go-grpc_out=internal/pb\
			--go-grpc_opt=paths=import\
			api/protobuf/*.proto

dependency-install:
	@go get $(GRPC_DEPS)
	@go install $(GRPC_DEPS)

modules:
	@git submodule update --init --remote
