
# to be run from the root of the project

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/messages.proto

# pip install grpcio-tools
python -m grpc_tools.protoc -I./proto --python_out=./client --pyi_out=./client --grpc_python_out=./client proto/messages.proto
