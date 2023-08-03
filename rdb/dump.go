package rdb

import (
	"bytes"

	"github.com/weedge/pkg/rdb/types"
)

// dump rdb entry obj
// https://github.com/redis/redis/blob/7.2-rc3/src/cluster.c#L6487
/* Write the footer, this is how it looks like:
 * ----------------+---------------------+---------------+
 * ... RDB payload | 2 bytes RDB version | 8 bytes CRC64 |
 * ----------------+---------------------+---------------+
 * RDB version and CRC are both in little endian.
 */

type String []byte
type List [][]byte
type Hash []struct {
	Field, Value []byte
}
type Set [][]byte
type ZSet []struct {
	Member []byte
	Score  float64
}

// RDB payload:
// | RDB TYPE STRING | len (encode string) | encode string |
func DumpStringValue(v String) []byte {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf)
	enc.EncodeType(types.RDBTypeString)
	enc.EncodeString(v)
	enc.EncodeDumpFooter()

	return buf.Bytes()
}

// RDB payload:
// | RDB TYPE HASH | len (encode hash)
// | len(encode field) | encode field | len(encode value) | encode value |.....|
func DumpHashValue(v Hash) []byte {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf)
	enc.EncodeType(types.RDBTypeHash)
	enc.EncodeLength(uint32(len(v)))
	for i := 0; i < len(v); i++ {
		enc.EncodeString(v[i].Field)
		enc.EncodeString(v[i].Value)
	}
	enc.EncodeDumpFooter()

	return buf.Bytes()
}

// RDB payload:
// | RDB TYPE LIST | len (encode list)
// | len(encode value) | encode value |.....|
func DumpListValue(v List) []byte {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf)
	enc.EncodeType(types.RDBTypeList)
	enc.EncodeLength(uint32(len(v)))
	for i := 0; i < len(v); i++ {
		enc.EncodeString(v[i])
	}
	enc.EncodeDumpFooter()

	return buf.Bytes()
}

// RDB payload:
// | RDB TYPE SET | len (encode set)
// | len(encode value) | encode value |.....|
func DumpSetValue(v Set) []byte {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf)
	enc.EncodeType(types.RDBTypeSet)
	enc.EncodeLength(uint32(len(v)))
	for i := 0; i < len(v); i++ {
		enc.EncodeString(v[i])
	}
	enc.EncodeDumpFooter()

	return buf.Bytes()
}

// RDB payload:
// | RDB TYPE ZSET | len (encode zset)
// | len(encode member) | encode member | len(encode score) | encode score |.....|
func DumpZSetValue(v ZSet) []byte {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf)
	enc.EncodeType(types.RDBTypeZSet)
	enc.EncodeLength(uint32(len(v)))
	for i := 0; i < len(v); i++ {
		enc.EncodeString(v[i].Member)
		enc.EncodeFloat(v[i].Score)
	}
	enc.EncodeDumpFooter()

	return buf.Bytes()
}
