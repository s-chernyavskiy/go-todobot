.PHONY: gen
gen:
	protoc -I api/protos api/protos/proto/todobot/*.proto \
	--go_out=./api/protos/gen/go \
	--go_opt=paths=source_relative \
	--go-grpc_out=./api/protos/gen/go \
	--go-grpc_opt=paths=source_relative

.PHONY: run
run:
	go run cmd/todobot/main.go \
	--config=./config/local.yaml