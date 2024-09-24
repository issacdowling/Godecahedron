package main

import (
	"encoding/json"
	"fmt"
	"net"

	comms "gitlab.com/issacdowling/godecahedron/internal/comms"
	types "gitlab.com/issacdowling/godecahedron/internal/types"
)

const (
	ADDR = "localhost:25565"
	// PORT = 25565
)

func main() {
	listener, err := net.Listen("tcp", ADDR)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Client connected: %s\n", conn.RemoteAddr().String())

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	msg := make([]byte, 2048)
	packet_len, err := conn.Read(msg)
	if err != nil {
		panic(err)
	}
	fmt.Println("Parsing claimed packet length:")
	claimed_length, length_bytes, _ := (types.NextVarint(msg[:packet_len]))
	fmt.Println(claimed_length)
	if claimed_length != int32(packet_len)-int32(length_bytes) {
		// panic("Claimed packet length was incorrect")
		fmt.Printf("Claimed packet length was incorrect: %v != %v", claimed_length, int32(packet_len)-int32(length_bytes))
	}

	comms.ParsePacket(msg[length_bytes:claimed_length])
	// conn.Close()
	//
	response, err := json.Marshal(types.StatusResponse{
		Version: types.SRVersion{
			Name:     "Test",
			Protocol: 767,
		},
		Players: types.SRPlayers{
			Max:    42,
			Online: 43,
			// No sample
		},
		Description: types.SRDesc{
			Text: "Godecahedron",
		},
		// No Favicon
		EnforcesSecureChat: false,
	})
	if err != nil {
		panic("Failed to JSON encode response: " + err.Error())
	}

	comms.SendPacket([]byte{0x00}, response, conn)
}
