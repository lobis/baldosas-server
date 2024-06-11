package main

import (
	pb "baldosas/proto"
	"baldosas/protocol"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"net"
	"sort"
	"sync"
	"time"
)

type position struct {
	x, y int
}

type baldosaServer struct {
	ipAddress   string
	connection  net.Conn
	stopChannel chan bool
}

var baldosas = map[position]*baldosaServer{
	{0, 0}: &baldosaServer{ipAddress: "192.168.31.43"},
	{1, 0}: &baldosaServer{ipAddress: "192.168.31.47"},
	{2, 0}: &baldosaServer{ipAddress: "192.168.31.51"},
	{3, 0}: &baldosaServer{ipAddress: "192.168.31.55"},
	{4, 0}: &baldosaServer{ipAddress: "192.168.31.59"},
	{5, 0}: &baldosaServer{ipAddress: "192.168.31.63"},
	{0, 1}: &baldosaServer{ipAddress: "192.168.31.42"},
	{1, 1}: &baldosaServer{ipAddress: "192.168.31.46"},
	{2, 1}: &baldosaServer{ipAddress: "192.168.31.50"},
	{3, 1}: &baldosaServer{ipAddress: "192.168.31.54"},
	{4, 1}: &baldosaServer{ipAddress: "192.168.31.58"},
	{5, 1}: &baldosaServer{ipAddress: "192.168.31.62"},
	{0, 2}: &baldosaServer{ipAddress: "192.168.31.41"},
	{1, 2}: &baldosaServer{ipAddress: "192.168.31.45"},
	{2, 2}: &baldosaServer{ipAddress: "192.168.31.49"},
	{3, 2}: &baldosaServer{ipAddress: "192.168.31.53"},
	{4, 2}: &baldosaServer{ipAddress: "192.168.31.57"},
	{5, 2}: &baldosaServer{ipAddress: "192.168.31.61"},
	{0, 3}: &baldosaServer{ipAddress: "192.168.31.40"},
	{1, 3}: &baldosaServer{ipAddress: "192.168.31.44"},
	{2, 3}: &baldosaServer{ipAddress: "192.168.31.48"},
	{3, 3}: &baldosaServer{ipAddress: "192.168.31.52"},
	{4, 3}: &baldosaServer{ipAddress: "192.168.31.56"},
	{5, 3}: &baldosaServer{ipAddress: "192.168.31.60"},
}

// map to relate position to string
var positionToString = map[position]string{
	{0, 3}: "A1", {1, 3}: "B1", {2, 3}: "C1", {3, 3}: "D1", {4, 3}: "E1", {5, 3}: "F1",
	{0, 2}: "A2", {1, 2}: "B2", {2, 2}: "C2", {3, 2}: "D2", {4, 2}: "E2", {5, 2}: "F2",
	{0, 1}: "A3", {1, 1}: "B3", {2, 1}: "C3", {3, 1}: "D3", {4, 1}: "E3", {5, 1}: "F3",
	{0, 0}: "A4", {1, 0}: "B4", {2, 0}: "C4", {3, 0}: "D4", {4, 0}: "E4", {5, 0}: "F4",
}

var baldosaPort = 51234

type grpcServer struct {
	pb.UnimplementedPositionsServiceServer
	pb.UnimplementedSensorServiceServer
	pb.UnimplementedLightServiceServer
	pb.UnimplementedSetLightsServiceServer
	pb.UnimplementedSetLightsStreamServiceServer
	pb.UnimplementedSetBrightnessServiceServer
}

func (s *grpcServer) GetPositions(_ context.Context, _ *pb.Empty) (*pb.Positions, error) {
	// iterate over the keys of baldosas to build positions
	positions := make([]*pb.Position, 0)
	lightsMutex.Lock()
	for pos := range lights {
		positions = append(positions, &pb.Position{X: int32(pos.x), Y: int32(pos.y)})
	}
	lightsMutex.Unlock()
	return &pb.Positions{Positions: positions}, nil
}

func setLightsHelper(in *pb.LightsStatus) {
	lightsMap := make(map[position]map[int]protocol.Light) // 3x3 to index to light
	lightsMutex.Lock()
	for _, light := range in.Lights {
		pos := position{x: int(light.Position.X), y: int(light.Position.Y)}
		positionSmall := positionBigToSmall(pos)
		index := positionToIndex(pos)

		if lightsMap[positionSmall] == nil {
			lightsMap[positionSmall] = make(map[int]protocol.Light)
		}
		lightsMap[positionSmall][index] = protocol.Light{
			Active: protocol.Color{
				R: uint8(light.Status.Active.R),
				G: uint8(light.Status.Active.G),
				B: uint8(light.Status.Active.B),
			},
			Inactive: protocol.Color{
				R: uint8(light.Status.Inactive.R),
				G: uint8(light.Status.Inactive.G),
				B: uint8(light.Status.Inactive.B),
			},
		}
	}
	for pos, lights3x3 := range lightsMap {
		baldosa := baldosas[pos]
		fmt.Println("Setting lights for", pos, "on", baldosa.ipAddress, "to", lights3x3)
		if baldosa.connection != nil {
			go func(baldosa baldosaServer, lights map[int]protocol.Light) {
				err := protocol.SendMessage(baldosa.connection, protocol.SetLights(lights))
				if err != nil {
					fmt.Println("Error sending message:", err)
				}
			}(*baldosa, lights3x3)
		}
	}
	lightsMutex.Unlock()
}

func (s *grpcServer) SetLights(_ context.Context, in *pb.LightsStatus) (*pb.Empty, error) {
	setLightsHelper(in)
	return &pb.Empty{}, nil
}

func (s *grpcServer) SetLightsStream(stream pb.SetLightsStreamService_SetLightsStreamServer) error {
	for {
		lightsStatus, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Println("Error receiving lights status:", err)
			return err
		}
		setLightsHelper(lightsStatus)
	}
}

func (s *grpcServer) SetBrightness(_ context.Context, in *pb.BrightnessStatus) (*pb.Empty, error) {
	for _, brightness := range in.Brightness {
		brightnessLevel := uint8(brightness.Brightness)
		pos := position{x: int(brightness.Position.X), y: int(brightness.Position.Y)}
		pos3x3 := positionBigToSmall(pos)
		fmt.Println("Setting brightness for", pos, "/", pos3x3, "to", brightness)
		// if position not in baldosas keys
		if baldosas[pos3x3] == nil {
			fmt.Println("Error: 3x3 position not found", pos3x3, "from", pos)
			continue
		}
		baldosa := baldosas[pos3x3]
		if baldosa.connection != nil {
			go func(b baldosaServer, brightness uint8) {
				err := protocol.SendMessage(b.connection, protocol.SetBrightness(brightness))
				if err != nil {
					fmt.Println("Error sending message:", err)
				}
			}(*baldosa, brightnessLevel)
		}
	}

	return &pb.Empty{}, nil
}

func (s *grpcServer) GetSensorStatusUpdates(_ *pb.Empty, stream pb.SensorService_GetSensorStatusUpdatesServer) error {
	for {
		select {
		case pos := <-sensorsUpdateChannel:
			sensorsMutex.Lock()
			err := stream.Send(&pb.SensorStatus{Position: &pb.Position{X: int32(pos.x), Y: int32(pos.y)}, Status: sensors[pos]})
			sensorsMutex.Unlock()
			if err != nil {
				fmt.Println("Error sending sensor status:", err)
				return err
			}
		}
	}
}

func (s *grpcServer) GetLightStatusUpdates(_ *pb.Empty, stream pb.LightService_GetLightStatusUpdatesServer) error {
	for {
		select {
		case pos := <-lightsUpdateChannel:
			lightsMutex.Lock()
			err := stream.Send(&pb.LightStatus{
				Position: &pb.Position{X: int32(pos.x), Y: int32(pos.y)}, //
				Status: &pb.Light{ //
					Active:   &pb.Color{R: int32(lights[pos].Active.R), G: int32(lights[pos].Active.G), B: int32(lights[pos].Active.B)},
					Inactive: &pb.Color{R: int32(lights[pos].Inactive.R), G: int32(lights[pos].Inactive.G), B: int32(lights[pos].Inactive.B)}}})
			lightsMutex.Unlock()
			if err != nil {
				fmt.Println("Error sending sensor status:", err)
				return err
			}
		}
	}
}

func startGrpcServer() {
	server := grpc.NewServer()
	pb.RegisterPositionsServiceServer(server, &grpcServer{})
	pb.RegisterSensorServiceServer(server, &grpcServer{})
	pb.RegisterLightServiceServer(server, &grpcServer{})
	pb.RegisterSetLightsServiceServer(server, &grpcServer{})
	pb.RegisterSetLightsStreamServiceServer(server, &grpcServer{})
	pb.RegisterSetBrightnessServiceServer(server, &grpcServer{})
	// start tcp server
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			fmt.Println("Error closing listener:", err)
		}
	}(listen)
	err = server.Serve(listen)
	if err != nil {
		fmt.Println("Error serving:", err)
		return
	}
}

var sensors = make(map[position]bool)
var sensorsMutex sync.Mutex
var sensorsUpdateChannel = make(chan position, 100)

var lights = make(map[position]protocol.Light)
var lightsMutex sync.Mutex
var lightsUpdateChannel = make(chan position, 100)

func readMessages(pos position, baldosa baldosaServer) {
	// read bytes one by one
	fmt.Println("Reading messages from", baldosa.ipAddress)
	for {
		select {
		case <-baldosa.stopChannel:
			return
		default:
			conn := baldosa.connection
			buf := make([]byte, 1)
			_, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading:", err)
				continue
			}
			if buf[0] != protocol.MessageBegin {
				continue
			}

			_, err = conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading:", err)
				continue
			}
			length := int(buf[0])

			_, err = conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading:", err)
				continue
			}
			messageType := int(buf[0])

			// read payload
			payload := make([]byte, length)
			_, err = conn.Read(payload)
			if err != nil {
				fmt.Println("Error reading:", err)
				continue
			}

			// read end of message
			_, err = conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading:", err)
				continue
			}
			if buf[0] != protocol.MessageEnd {
				fmt.Println("Error: expected end of message")
				continue
			}

			// process payload
			switch messageType {
			case protocol.MessageTypePong:
				// fmt.Println("Received pong")
			case protocol.MessageTypeSensorsStatus:
				// fmt.Println("Received sensors status")
				sensorsMutex.Lock()
				// entries are length of payload / 2
				entries := len(payload) / 2
				for i := 0; i < entries; i++ {
					index := payload[i*2]
					value := payload[i*2+1]
					previousValue := sensors[indexToPosition(int(index), pos)]
					sensors[indexToPosition(int(index), pos)] = value == 1
					if previousValue != (value == 1) {
						fmt.Println("Sensor", indexToPosition(int(index), pos), "changed to", value == 1)
						select {
						case sensorsUpdateChannel <- indexToPosition(int(index), pos):
						default:
						}
					}
				}
				sensorsMutex.Unlock()
			case protocol.MessageTypeLightsStatus:
				fmt.Println("Received lights status")
				lightsMutex.Lock()
				entries := len(payload) / 7
				for i := 0; i < entries; i++ {
					index := payload[i*7]
					active := protocol.Color{
						R: payload[i*7+1],
						G: payload[i*7+2],
						B: payload[i*7+3],
					}
					inactive := protocol.Color{
						R: payload[i*7+4],
						G: payload[i*7+5],
						B: payload[i*7+6],
					}
					previousValue := lights[indexToPosition(int(index), pos)]
					lights[indexToPosition(int(index), pos)] = protocol.Light{
						Active:   active,
						Inactive: inactive,
					}
					if previousValue != lights[indexToPosition(int(index), pos)] {
						fmt.Println("Light", indexToPosition(int(index), pos), "changed to", lights[indexToPosition(int(index), pos)])
						select {
						case lightsUpdateChannel <- indexToPosition(int(index), pos):
						default:
						}
					}
				}
				lightsMutex.Unlock()
			default:
				fmt.Println("Error: unknown message type:", messageType)
			}
		}
	}
}

// Converts an index in a 3x3 grid to a global position considering the position of the 3x3 grid with horizontal flip
func indexToPosition(index int, positionOf3x3 position) position {
	return position{
		x: positionOf3x3.x*3 + (2 - index%3),
		y: positionOf3x3.y*3 + (2 - index/3),
	}
}

// Converts a global position to a 3x3 grid position
func positionBigToSmall(pos position) position {
	return position{
		x: pos.x / 3,
		y: pos.y / 3,
	}
}

// Converts a position within a 3x3 grid to an index considering horizontal flip
func positionToIndex(pos position) int {
	localX := pos.x % 3
	localY := pos.y % 3
	return (2-localY)*3 + (2 - localX)
}

func main() {
	go startGrpcServer()

	// initialize sensors and lights
	for key := range baldosas {
		for i := 0; i < 9; i++ {
			sensors[indexToPosition(i, key)] = false
			lights[indexToPosition(i, key)] = protocol.Light{
				Active:   protocol.Color{R: 255, G: 255, B: 255}, // white
				Inactive: protocol.Color{R: 0, G: 0, B: 255},     // blue
			}
		}
	}

	for pos, _ := range baldosas {
		// mutable ref to baldosa
		baldosas[pos] = &baldosaServer{
			ipAddress:   baldosas[pos].ipAddress,
			connection:  nil,
			stopChannel: make(chan bool),
		}

		go func(pos position) {
			for {
				if baldosas[pos].connection != nil {
					// fmt.Println("Sending ping to", baldosas[pos].ipAddress)
					err := protocol.SendMessage(baldosas[pos].connection, protocol.Ping())
					if err != nil {
						fmt.Println("Error sending ping:", err)
						baldosas[pos].connection = nil
						baldosas[pos].stopChannel <- true
					}
					time.Sleep(1 * time.Second)
				} else {
					fmt.Println("Connecting to", baldosas[pos].ipAddress, "on port", baldosaPort)
					conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", baldosas[pos].ipAddress, baldosaPort))
					if err != nil {
						fmt.Println("Error establishing connection:", err)
						baldosas[pos].connection = nil
						if baldosas[pos].stopChannel != nil {
							close(baldosas[pos].stopChannel)
							baldosas[pos].stopChannel = nil
						}
					} else {
						fmt.Println("Connected to", baldosas[pos].ipAddress, "on port", baldosaPort)
						baldosas[pos].connection = conn
						baldosas[pos].stopChannel = make(chan bool)

						go readMessages(pos, *baldosas[pos])

						err := protocol.SendMessage(baldosas[pos].connection, protocol.RequestSensorsStatus())
						if err != nil {
							fmt.Println("Error sending message:", err)
						}

						err = protocol.SendMessage(baldosas[pos].connection, protocol.RequestLightsStatus())
						if err != nil {
							fmt.Println("Error sending message:", err)
						}
					}
				}
			}
		}(pos)
	}

	// create a goroutine that prints how many connections
	go func() {
		for {
			fmt.Println("Connections summary:")
			count := 0
			for _, baldosa := range baldosas {
				if baldosa.connection != nil {
					count++
				}
			}
			// print the number of connections out of total
			fmt.Println("Connected to", count, "out of", len(baldosas))
			// print connected ips, and not connected ips
			// sets of connected strings, disconnected
			var connected []string
			var disconnected []string

			for pos, baldosa := range baldosas {
				if baldosa.connection != nil {
					connected = append(connected, positionToString[pos])
				} else {
					disconnected = append(disconnected, positionToString[pos])
				}
			}

			// sort both alphabetically
			sort.Strings(connected)
			sort.Strings(disconnected)

			fmt.Println("Connected to:", connected)
			fmt.Println("Not connected to:", disconnected)

			time.Sleep(5 * time.Second)
		}
	}()

	select {} // Block forever
}
