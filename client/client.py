import grpc

import messages_pb2 as messages
import messages_pb2_grpc as messages_grpc


def main():
    # connect to server
    channel = grpc.insecure_channel('localhost:50051')

    # connect to the sensor status stream
    sensor_status_stream = messages_grpc.SensorServiceStub(channel)
    sensor_status = sensor_status_stream.GetSensorStatusUpdates(messages.Empty())

    # print the sensor status
    for status in sensor_status:
        print(status)


if __name__ == "__main__":
    main()
