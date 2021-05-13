.PHONY: grpc-gen
grpc-gen:
	protoc --go_out=. --go-grpc_out=. proto/health.proto
	python3 -m grpc_tools.protoc -Iproto --python_out=./proto --grpc_python_out=./proto proto/health.proto
