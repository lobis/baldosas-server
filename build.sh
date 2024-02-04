
# to be run from the root of the project

protoc --go_out=./proto --go_opt=paths=source_relative --go-grpc_out=./proto --go-grpc_opt=paths=source_relative --js_out=import_style=commonjs,binary:./app/src/proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:./app/src/proto --proto_path=proto proto/messages.proto

# pip install grpcio-tools
python -m grpc_tools.protoc -I./proto --python_out=./client --pyi_out=./client --grpc_python_out=./client proto/messages.proto
