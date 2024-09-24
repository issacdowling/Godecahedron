package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
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
	claimed_length, length_bytes, _ := (nextVarint(msg[:packet_len]))
	fmt.Println(claimed_length)
	if claimed_length != int32(packet_len)-int32(length_bytes) {
		// panic("Claimed packet length was incorrect")
		fmt.Printf("Claimed packet length was incorrect: %v != %v", claimed_length, int32(packet_len)-int32(length_bytes))
	}

	parsePacket(msg[length_bytes:claimed_length])
	// conn.Close()
	//
	response, err := json.Marshal(StatusResponse{
		Version: SRVersion{
			Name:     "Test",
			Protocol: 767,
		},
		Players: SRPlayers{
			Max:    42,
			Online: 43,
			// No sample
		},
		Description: SRDesc{
			Text: "Godecahedron",
		},
		// No Favicon
		EnforcesSecureChat: false,
	})
	if err != nil {
		panic("Failed to JSON encode response: " + err.Error())
	}

	sendPacket(response, conn)
}

// Parse Minecraft VarInts.
// Needs to be assured that the first available byte is part of the varint, but will automatically
// stop once the end is reached, and return the number of bytes read. Errors if too long (over 5 bytes).
func nextVarint(buf []byte) (int32, int8, error) {
	const valueBits int8 = 7
	const maxLength int8 = valueBits * 4
	const valueBitmask byte = 0x7F
	const endBitmask byte = 0x80

	var value int32 = 0
	var pos int8 = 0
	for _, byte := range buf {
		value |= int32(byte&valueBitmask) << pos
		pos += valueBits

		if byte&endBitmask == 0 {
			break
		}

		if pos > maxLength {
			return 0, 0, errors.New("Varint may not be longer than 5 bytes")
		}
	}
	return value, pos / valueBits, nil
}

func writeVarint(value int32) ([]byte, int8) {
	var bytes int8
	var varint []byte
	const valueBitmask byte = 0x7F
	const endBitmask byte = 0x80
	for {
		bytes++
		fmt.Println(varint)
		if (byte(value) & ^valueBitmask) == 0 {
			varint = append(varint, byte(value))
			return varint, bytes
		}
		fmt.Println("\n\n\n\nXXXXX")

		varint = append(varint, byte(value)&valueBitmask|endBitmask)
		value >>= 7
	}
}

func parsePacket(buf []byte) {
	fmt.Println("Parsing packet ID:")
	packetId, length, err := nextVarint(buf)
	if err != nil {
		panic(err)
	}

	fmt.Println(packetId)

	fmt.Println("Parsing protocol version:")
	protoVer, length2, err := nextVarint(buf[length:])
	if err != nil {
		panic(err)
	}
	fmt.Println(protoVer, length2)

	fmt.Println("Parsing address length:")
	adrLen, length3, err := nextVarint(buf[length+length2:])
	if err != nil {
		panic(err)
	}
	fmt.Println(adrLen, length3)

	fmt.Println("Printing address:")
	fmt.Println(string(buf[length+length2+length3 : length+length2+length3+int8(adrLen)]))

	fmt.Println("Printing port:")
	portReader := bytes.NewReader(buf[length+length2+length3+int8(adrLen) : length+length2+length3+int8(adrLen)+2])
	var portInt uint16
	binary.Read(portReader, binary.BigEndian, &portInt)
	fmt.Println(portInt)
}

func sendPacket(buf []byte, conn net.Conn) {
	lenBufBytes, lenBufLen := writeVarint(int32(len(buf)))
	lenTotalBytes, _ := writeVarint(1 + int32(lenBufLen) + int32(len(buf)))
	var testResponse []byte

	testResponse = append(testResponse, lenTotalBytes...)
	testResponse = append(testResponse, byte(0x00))
	testResponse = append(testResponse, lenBufBytes...)
	testResponse = append(testResponse, buf...)
	fmt.Println(testResponse)
	conn.Write(testResponse)
}
