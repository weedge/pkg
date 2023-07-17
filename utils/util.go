package utils

import (
	"reflect"
	"regexp"
	"strings"
)

func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

func BuildMatchRegexp(match string) (*regexp.Regexp, error) {
	var err error
	var r *regexp.Regexp

	if len(match) > 0 {
		if r, err = regexp.Compile(match); err != nil {
			return nil, err
		}
	}

	return r, nil
}

// AsyncNoBlockSend async no block send notify channel.
func AsyncNoBlockSend(ch chan<- struct{}) {
	select {
	case ch <- struct{}{}:
	default:
	}
}

// AsyncNoBlockSendUint64 async no block send Uint64  channel.
func AsyncNoBlockSendUint64(ch chan uint64, v uint64) {
	select {
	case ch <- v:
	default:
	}
}

// MinInt
func MinInt(a int, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

// GetProto get proto by addr
func GetProto(addr string) string {
	if strings.Contains(addr, "/") {
		return "unix"
	} else {
		return "tcp"
	}
}
