# https://grpc.io/docs/languages/go/quickstart/
# brew install protoc-gen-grpc-web
# to be run from the root of the project

protoc --go_out=./proto --go_opt=paths=source_relative --go-grpc_out=./proto --go-grpc_opt=paths=source_relative --proto_path=proto proto/messages.proto

# pip install grpcio-tools
python -m grpc_tools.protoc -I./proto --python_out=./client --pyi_out=./client --grpc_python_out=./client proto/messages.proto
