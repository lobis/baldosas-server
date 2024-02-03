package main

import (
	pb "baldosas/proto"
	"baldosas/protocol"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
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

var baldosas = make(map[position]baldosaServer)

var baldosaPort = 1234

type grpcServer struct {
	pb.UnimplementedStatusServiceServer
}

func (s *grpcServer) GetConnectedClients(ctx context.Context, req *pb.Empty) (*pb.Status, error) {
	return &pb.Status{ConnectedClients: 10}, nil
}

func startGrpcServer() {
	server := grpc.NewServer()
	pb.RegisterStatusServiceServer(server, &grpcServer{})
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

var lights = make(map[position]protocol.Light)
var lightsMutex sync.Mutex

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
				fmt.Println("Received pong")
			case protocol.MessageTypeSensorsStatus:
				fmt.Println("Received sensors status")
				sensorsMutex.Lock()
				// entries are length of payload / 2
				entries := len(payload) / 2
				for i := 0; i < entries; i++ {
					index := payload[i*2]
					value := payload[i*2+1]
					sensors[indexToPosition(int(index), pos)] = value == 1
				}
				sensorsMutex.Unlock()
			case protocol.MessageTypeLightsStatus:
				fmt.Println("Received lights status")
				lightsMutex.Lock()
				entries := len(payload) / 7
				for i := 0; i < entries; i++ {
					index := payload[i*7]
					off := protocol.Color{
						R: payload[i*7+1],
						G: payload[i*7+2],
						B: payload[i*7+3],
					}
					on := protocol.Color{
						R: payload[i*7+4],
						G: payload[i*7+5],
						B: payload[i*7+6],
					}
					lights[indexToPosition(int(index), pos)] = protocol.Light{
						On:  on,
						Off: off,
					}
				}
				lightsMutex.Unlock()
			default:
				fmt.Println("Error: unknown message type:", messageType)
			}
		}
	}
}

func indexToPosition(index int, positionOf3x3 position) position {
	return position{
		x: positionOf3x3.x*3 + index%3,
		y: positionOf3x3.y*3 + index/3,
	}
}

func main() {
	go startGrpcServer()

	// initialize baldosas
	baldosas[position{x: 0, y: 0}] = baldosaServer{ipAddress: "192.168.1.139"}

	// initialize sensors and lights
	for key := range baldosas {
		for i := 0; i < 9; i++ {
			sensors[indexToPosition(i, key)] = false
			lights[indexToPosition(i, key)] = protocol.Light{
				On:  protocol.Color{R: 255, G: 255, B: 255},
				Off: protocol.Color{R: 0, G: 0, B: 0},
			}
		}
	}

	for pos, baldosa := range baldosas {
		baldosas[pos] = baldosaServer{
			ipAddress:   baldosa.ipAddress,
			connection:  nil,
			stopChannel: make(chan bool),
		}

		go func(pos position, baldosa baldosaServer) {
			for {
				if baldosa.connection != nil {
					err := protocol.SendMessage(baldosa.connection, protocol.Ping())
					if err != nil {
						fmt.Println("Error sending ping:", err)
						baldosa.connection = nil
						baldosa.stopChannel <- true
					}
					time.Sleep(1 * time.Second)
				} else {
					fmt.Println("Connecting to", baldosa.ipAddress, "on port", baldosaPort)
					conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", baldosa.ipAddress, baldosaPort))
					if err != nil {
						fmt.Println("Error establishing connection:", err)
						baldosa.connection = nil
						if baldosa.stopChannel != nil {
							close(baldosa.stopChannel)
							baldosa.stopChannel = nil
						}
					} else {
						fmt.Println("Connected to", baldosa.ipAddress, "on port", baldosaPort)
						baldosa.connection = conn
						baldosa.stopChannel = make(chan bool)

						go readMessages(pos, baldosa)

						err := protocol.SendMessage(baldosa.connection, protocol.RequestSensorsStatus())
						if err != nil {
							fmt.Println("Error sending message:", err)
						}

						err = protocol.SendMessage(baldosa.connection, protocol.RequestLightsStatus())
						if err != nil {
							fmt.Println("Error sending message:", err)
						}
					}
				}
			}
		}(pos, baldosa)
	}

	select {} // Block forever
}
