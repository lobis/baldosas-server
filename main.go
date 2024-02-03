package main

import (
	pb "baldosas/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"time"
)

type Server struct {
	address string
	port    int
}

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

func main() {
	// start grpc server
	go startGrpcServer()

	serverList := []Server{
		{address: "192.168.1.138", port: 1234},
	}

	connections := make([]net.Conn, len(serverList))
	for i, server := range serverList {
		go func(i int, server Server) {
			for {
				if connections[i] != nil {
					err := sendPing(connections[i])
					if err != nil {
						fmt.Println("Error sending ping:", err)
						connections[i] = nil
					}
					time.Sleep(1 * time.Second)
				}
				fmt.Println("Connecting to server", server.address, "on port", server.port)
				conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server.address, server.port))
				if err != nil {
					fmt.Println("Error connecting to server:", err)
					connections[i] = nil
				} else {
					connections[i] = conn
				}
			}
		}(i, server)
	}

	select {} // Block forever
}

func sendPing(conn net.Conn) error {
	_, err := conn.Write([]byte{0x0B, 0x00, 0x0C, 0x0A})
	return err
}
