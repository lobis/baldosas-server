import messages_pb2 as messages
import grpc
import messages_pb2_grpc as messages_grpc

def main():
    # connect to server
    channel = grpc.insecure_channel('localhost:50051')

    # create a stub
    stub = messages_grpc.StatusServiceStub(channel)

    # create a request
    request = messages.Empty()
    response = stub.GetConnectedClients(request)
    print(response.connectedClients)

if __name__ == "__main__":
    main()