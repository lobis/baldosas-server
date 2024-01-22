package main

import (
	"fmt"
	"net"
	"time"
)

type Server struct {
	address string
	port    int
}

func main() {
	fmt.Println("Starting client...")

	serverList := []Server{{address: "192.168.1.177", port: 1234}}
	connections := make([]net.Conn, len(serverList))
	for i, server := range serverList {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server.address, server.port))
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			return
		}
		connections[i] = conn
	}

	fmt.Println("Connected to server")

	for _, conn := range connections {
		// Start a goroutine for sending ping messages periodically
		go func(conn net.Conn) {
			defer func() {
				err := conn.Close()
				if err != nil {
					fmt.Println("Error closing connection:", err)
				}
			}()

			for {
				err := sendPing(conn)
				if err != nil {
					fmt.Println("Error sending ping:", err)
					break // Break the loop if there's an error
				}
				time.Sleep(1 * time.Second)
			}
		}(conn)

		// start goroutine for receiving messages
		go func(conn net.Conn) {
			defer func() {
				err := conn.Close()
				if err != nil {
					fmt.Println("Error closing connection:", err)
				}
			}()

			for {
				// read from the connection
				buf := make([]byte, 1)
				_, err := conn.Read(buf)
				if err != nil {
					fmt.Println("Error reading from connection:", err)
					break
				}
				fmt.Println("Received:", buf)
			}
		}(conn)
	}

	select {} // Block forever
}

func sendPing(conn net.Conn) error {
	// Modify this function to send your "ping" message
	fmt.Println("Sending ping")
	// send two 0x0B bytes
	_, err := conn.Write([]byte{0x0B, 0x00, 0x0C, 0x0A})
	return err
}
