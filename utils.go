package gsstat

import (
	"encoding/binary"
	"math"
	"strconv"
)

// BytesToUint16 - expect 2 bytes
func BytesToUint16(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)
}

// BytesToUint32 - expect 4 bytes
func BytesToUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

// BytesToUint64 - expect 8 bytes
func BytesToUint64(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}

// BytesToFloat32 - expect 4 bytes
func BytesToFloat32(b []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(b))
}

// StringToUint8 ..
func StringToUint8(s string) uint8 {
	tInt64, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return 0
	}
	return uint8(tInt64)
}

// StringToUint16 ..
func StringToUint16(s string) uint16 {
	tInt64, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return 0
	}
	return uint16(tInt64)
}

// ByteToBool ..
func ByteToBool(b byte) bool {
	return b != 0
}
