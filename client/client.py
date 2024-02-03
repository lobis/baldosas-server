import grpc

import messages_pb2 as messages
import messages_pb2_grpc as messages_grpc


def main():
    # connect to server
    channel = grpc.insecure_channel('localhost:50051')

    # positions
    positions_stub = messages_grpc.PositionsServiceStub(channel)
    positions = positions_stub.GetPositions(messages.Empty())
    print(positions)

    # send light update
    light_stub = messages_grpc.SetLightsServiceStub(channel)
    light_status = messages.LightStatus(position=messages.Position(x=1, y=1),
                                        status=messages.Light(on=messages.Color(r=255, g=0, b=0),
                                                              off=messages.Color(r=255, g=0, b=255)))
    light_stub.SetLights(messages.LightsStatus(lights=[light_status]))

    # connect to the sensor status stream
    sensor_status_stream = messages_grpc.SensorServiceStub(channel)
    sensor_status = sensor_status_stream.GetSensorStatusUpdates(messages.Empty())

    # light status
    light_status_stream = messages_grpc.LightServiceStub(channel)
    light_status = light_status_stream.GetLightStatusUpdates(messages.Empty())

    # print the light status
    for status in light_status:
        print(status)


if __name__ == "__main__":
    main()
