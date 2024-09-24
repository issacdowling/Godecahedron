package types

import "errors"

// Tried using anonymous structs, but declaring them requires restating their types
type StatusResponse struct {
	Version            SRVersion `json:"version"`
	Players            SRPlayers `json:"players"`
	Description        SRDesc    `json:"description"`
	Favicon            string    `json:"favicon"`
	EnforcesSecureChat bool      `json:"enforcesSecureChat"`
}

type SRVersion struct {
	Name     string `json:"name"`
	Protocol int16  `json:"protocol"`
}

type SRPlayers struct {
	Max    int8 `json:"max"`
	Online int8 `json:"online"`
	Sample []struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	} `json:"sample"`
}

type SRDesc struct {
	Text string `json:"text"`
}

// Parse Minecraft VarInts.
// Needs to be assured that the first available byte is part of the varint, but will automatically
// stop once the end is reached, and return the number of bytes read. Errors if too long (over 5 bytes).
func NextVarint(buf []byte) (int32, int8, error) {
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

func WriteVarint(value int32) ([]byte, int8) {
	var bytes int8
	var varint []byte
	const valueBitmask byte = 0x7F
	const endBitmask byte = 0x80
	for {
		bytes++
		if (byte(value) & ^valueBitmask) == 0 {
			varint = append(varint, byte(value))
			return varint, bytes
		}

		varint = append(varint, byte(value)&valueBitmask|endBitmask)
		value >>= 7
	}
}
