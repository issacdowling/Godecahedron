package comms

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	types "gitlab.com/issacdowling/godecahedron/internal/types"
)

func ParsePacket(buf []byte) {
	fmt.Println("Parsing packet ID:")
	packetId, length, err := types.NextVarint(buf)
	if err != nil {
		panic(err)
	}

	fmt.Println(packetId)

	fmt.Println("Parsing protocol version:")
	protoVer, length2, err := types.NextVarint(buf[length:])
	if err != nil {
		panic(err)
	}
	fmt.Println(protoVer, length2)

	fmt.Println("Parsing address length:")
	adrLen, length3, err := types.NextVarint(buf[length+length2:])
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

func SendPacket(id []byte, data []byte, conn net.Conn) {
	lenBufBytes, lenBufLen := types.WriteVarint(int32(len(data)))
	lenTotalBytes, _ := types.WriteVarint(1 + int32(lenBufLen) + int32(len(data)))
	var testResponse []byte

	testResponse = append(testResponse, lenTotalBytes...)
	testResponse = append(testResponse, byte(0x00))
	testResponse = append(testResponse, lenBufBytes...)
	testResponse = append(testResponse, data...)
	conn.Write(testResponse)
}
