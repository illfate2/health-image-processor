.PHONY: grpc-gen
grpc-gen:
	protoc --go_out=. --go-grpc_out=. --python_out=. --grpc_python_out=. proto/health.proto
