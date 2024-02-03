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

type baldosa struct {
	ipAddress string
	position  position
}

var baldosas = []baldosa{
	{
		ipAddress: "192.168.1.139",
		position:  position{x: 0, y: 0},
	},
	{
		ipAddress: "192.168.1.138",
		position:  position{x: 1, y: 0},
	},
}

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

var sensorsMutex sync.Mutex

func readMessages(conn net.Conn, stop chan bool) {
	// read bytes one by one
	for {
		select {
		case <-stop:
			return
		default:
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
				fmt.Println("Received sensors status")
				sensorsMutex.Lock()
				// entries are length of payload / 2
				entries := len(payload) / 2
				for i := 0; i < entries; i++ {
					index := payload[i*2]
					value := payload[i*2+1]
					fmt.Println("Sensor", index, "value", value)
				}
				sensorsMutex.Unlock()
			default:
				fmt.Println("Error: unknown message type:", payload[0])
			}
		}
	}
}

func main() {
	sensors := make(map[position]bool)
	for _, entry := range baldosas {
		initialPosition := position{
			x: entry.position.x * 3,
			y: entry.position.y * 3,
		}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				position := position{x: initialPosition.x + i, y: initialPosition.y + j}
				if _, ok := sensors[position]; ok {
					fmt.Println("Error: repeated position", position)
					return
				}
				sensors[position] = false
			}
		}
	}
	// print sensors
	for key, value := range sensors {
		fmt.Println("Sensor", key, "value", value)
	}
	// start grpc server
	go startGrpcServer()

	connections := make([]net.Conn, len(baldosas))
	stopChannels := make([]chan bool, len(baldosas))
	for i, server := range baldosas {
		stopChannels[i] = make(chan bool)
		go func(i int, this baldosa) {
			for {
				if connections[i] != nil {
					err := sendPing(connections[i])
					if err != nil {
						fmt.Println("Error sending ping:", err)
						connections[i] = nil
						stopChannels[i] <- true
					}
					time.Sleep(1 * time.Second)
				} else {
					fmt.Println("Connecting to", this.ipAddress, "on port", baldosaPort)
					conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", this.ipAddress, baldosaPort))
					if err != nil {
						fmt.Println("Error establishing connection:", err)
						connections[i] = nil
						if stopChannels[i] != nil {
							close(stopChannels[i])
							stopChannels[i] = nil
						}
					} else {
						fmt.Println("Connected to", this.ipAddress, "on port", baldosaPort)
						connections[i] = conn
						stopChannels[i] = make(chan bool)
						go readMessages(conn, stopChannels[i])
					}
				}
			}
		}(i, server)
	}

	select {} // Block forever
}

func sendPing(conn net.Conn) error {
	_, err := conn.Write(protocol.Ping())
	return err
}
