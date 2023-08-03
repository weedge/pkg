package structure

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/weedge/pkg/utils/logutils"
)

const (
	RDB6BitLen      = 0 // RDB_6BITLEN
	RDB14BitLen     = 1 // RDB_14BITLEN
	RDB32or64BitLen = 2
	RDBEncVal       = 3 // RDB_ENCVAL
	RDB32BitLen     = 0x80
	RDB64BitLen     = 0x81
)

func ReadLength(rd io.Reader) uint64 {
	length, special, err := readEncodedLength(rd)
	if special {
		logutils.Criticalf("illegal length special=true, encoding: %d", length)
	}
	if err != nil {
		logutils.CriticalError(err)
	}
	return length
}

func readEncodedLength(rd io.Reader) (length uint64, special bool, err error) {
	var lengthBuffer = make([]byte, 8)

	firstByte := ReadByte(rd)
	first2bits := (firstByte & 0xc0) >> 6 // first 2 bits of encoding
	switch first2bits {
	case RDB6BitLen:
		length = uint64(firstByte) & 0x3f
	case RDB14BitLen:
		nextByte := ReadByte(rd)
		length = (uint64(firstByte)&0x3f)<<8 | uint64(nextByte)
	case RDB32or64BitLen:
		if firstByte == RDB32BitLen {
			_, err = io.ReadFull(rd, lengthBuffer[0:4])
			if err != nil {
				return 0, false, fmt.Errorf("read len32Bit failed: %s", err.Error())
			}
			length = uint64(binary.BigEndian.Uint32(lengthBuffer))
		} else if firstByte == RDB64BitLen {
			_, err = io.ReadFull(rd, lengthBuffer)
			if err != nil {
				return 0, false, fmt.Errorf("read len64Bit failed: %s", err.Error())
			}
			length = binary.BigEndian.Uint64(lengthBuffer)
		} else {
			return 0, false, fmt.Errorf("illegal length encoding: %x", firstByte)
		}
	case RDBEncVal:
		special = true
		length = uint64(firstByte) & 0x3f
	}
	return length, special, nil
}
