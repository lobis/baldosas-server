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
                                                              off=messages.Color(r=255, g=0, b=252)))
    light_stub.SetLights(messages.LightsStatus(lights=[light_status]))

    # send a stream of light status to the stream service
    stream_stub = messages_grpc.SetLightsStreamServiceStub(channel)
    light_status = messages.LightStatus(position=messages.Position(x=1, y=1),
                                        status=messages.Light(on=messages.Color(r=255, g=0, b=0),
                                                              off=messages.Color(r=255, g=0, b=252)))

    def generate_light_status():
        for i in range(10):
            yield messages.LightsStatus(lights=[light_status])

    stream = stream_stub.SetLightsStream(generate_light_status())


if __name__ == "__main__":
    main()
