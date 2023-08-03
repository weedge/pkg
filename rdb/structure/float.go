package structure

import (
	"encoding/binary"
	"io"
	"math"
	"strconv"

	"github.com/weedge/pkg/utils/logutils"
)

func ReadFloat(rd io.Reader) float64 {
	u := ReadUint8(rd)

	switch u {
	case 253:
		return math.NaN()
	case 254:
		return math.Inf(0)
	case 255:
		return math.Inf(-1)
	default:
		buf := make([]byte, u)
		_, err := io.ReadFull(rd, buf)
		if err != nil {
			return 0
		}

		v, err := strconv.ParseFloat(string(buf), 64)
		if err != nil {
			logutils.CriticalError(err)
		}
		return v
	}
}

func ReadDouble(rd io.Reader) float64 {
	var buf = make([]byte, 8)
	_, err := io.ReadFull(rd, buf)
	if err != nil {
		logutils.CriticalError(err)
	}
	num := binary.LittleEndian.Uint64(buf)
	return math.Float64frombits(num)
}