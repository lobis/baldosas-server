import grpc
import random
import threading

import messages_pb2 as messages
import messages_pb2_grpc as messages_grpc

state = dict()


def update_state(channel: grpc.Channel):
    global state
    # read from the sensors update channel
    sensor_stub = messages_grpc.SensorServiceStub(channel)
    # read from stream forever
    lights_stub = messages_grpc.SetLightsServiceStub(channel)
    brightness_stub = messages_grpc.SetBrightnessServiceStub(channel)

    # count how many positions for each color
    color_count = {
        (color.r, color.g, color.b): 0 for color in state.values()
    }
    for color in state.values():
        color_count[(color.r, color.g, color.b)] += 1

    original_positions = list(state.keys())
    for sensor in sensor_stub.GetSensorStatusUpdates(messages.Empty()):
        position = (sensor.position.x, sensor.position.y)
        new_state = sensor.status
        if not new_state:
            continue
        if position not in state:
            continue
        # remove the entry with key "position"
        print("Killing light at", position)
        color = state[position]
        color_count[(color.r, color.g, color.b)] -= 1
        state.pop(position, None)
        lights_message = messages.LightsStatus(lights=[
            messages.LightStatus(position=messages.Position(x=position[0], y=position[1]),
                                 status=messages.Light(inactive=messages.Color(r=0, g=0, b=0),
                                                       active=messages.Color(r=0, g=0, b=0)))
        ])
        lights_stub.SetLights(lights_message)

        if 0 in color_count.values():
            break

    # game finished
    print("Game finished")
    # set all to the color of the winner
    color = min(color_count, key=color_count.get)
    print("Winner color", color)
    lights_message = messages.LightsStatus(lights=[
        messages.LightStatus(position=messages.Position(x=position[0], y=position[1]),
                             status=messages.Light(active=messages.Color(r=0, g=0, b=0),
                                                   inactive=messages.Color(r=color[0], g=color[1], b=color[2])))
        for position in original_positions
    ])
    lights_stub.SetLights(lights_message)


def main():
    # connect to server
    # host = "192.168.31.187"
    host = "localhost"
    channel = grpc.insecure_channel(f"{host}:50051")
    print("Connected to", host)

    # positions
    positions_stub = messages_grpc.PositionsServiceStub(channel)
    positions = positions_stub.GetPositions(messages.Empty())
    # print(positions)

    colors = [
        # orange
        messages.Color(r=255, g=127, b=0),
        # cyan
        messages.Color(r=0, g=255, b=255),
        # magenta
        messages.Color(r=255, g=0, b=255),
    ]

    sequence = colors * (len(positions.positions) // len(colors) + 1)
    # remove from the back until we have exactly the number of positions
    sequence = sequence[:len(positions.positions)]
    # random shuffle
    random.shuffle(sequence)
    global state

    state = {
        (position.x, position.y): color for position, color in zip(positions.positions, sequence)
    }

    # set brightness
    brightness = 100
    brightness_stub = messages_grpc.SetBrightnessServiceStub(channel)
    brightness_message = messages.BrightnessStatus(brightness=[
        messages.Brightness(position=messages.Position(x=position[0], y=position[1]),
                             brightness=brightness)
        for position, color in state.items()
    ])
    brightness_stub.SetBrightness(brightness_message)

    lights_stub = messages_grpc.SetLightsServiceStub(channel)
    lights_message = messages.LightsStatus(lights=[
        messages.LightStatus(position=messages.Position(x=position[0], y=position[1]),
                             status=messages.Light(inactive=color, active=messages.Color(r=0, g=0, b=0)))
        for position, color in state.items()
    ])
    lights_stub.SetLights(lights_message)

    # start a thread to update state
    update_state_thread = threading.Thread(target=update_state, args=(channel,), daemon=True)
    update_state_thread.start()

    while True:
        pass


if __name__ == "__main__":
    main()
