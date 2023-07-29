package rdb

import (
	"bytes"

	"github.com/cupcake/rdb"
)

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

func DumpStringValue(v String) []byte {
	buf := bytes.NewBuffer(nil)
	enc := rdb.NewEncoder(buf)
	enc.EncodeType(rdb.TypeString)
	enc.EncodeString(v)
	enc.EncodeDumpFooter()

	return buf.Bytes()
}

func DumpHashValue(v Hash) []byte {
	buf := bytes.NewBuffer(nil)
	enc := rdb.NewEncoder(buf)
	enc.EncodeType(rdb.TypeHash)
	enc.EncodeLength(uint32(len(v)))
	for i := 0; i < len(v); i++ {
		enc.EncodeString(v[i].Field)
		enc.EncodeString(v[i].Value)
	}
	enc.EncodeDumpFooter()

	return buf.Bytes()
}

func DumpListValue(v List) []byte {
	buf := bytes.NewBuffer(nil)
	enc := rdb.NewEncoder(buf)
	enc.EncodeType(rdb.TypeList)
	enc.EncodeLength(uint32(len(v)))
	for i := 0; i < len(v); i++ {
		enc.EncodeString(v[i])
	}
	enc.EncodeDumpFooter()

	return buf.Bytes()
}

func DumpSetValue(v Set) []byte {
	buf := bytes.NewBuffer(nil)
	enc := rdb.NewEncoder(buf)
	enc.EncodeType(rdb.TypeSet)
	enc.EncodeLength(uint32(len(v)))
	for i := 0; i < len(v); i++ {
		enc.EncodeString(v[i])
	}
	enc.EncodeDumpFooter()

	return buf.Bytes()
}

func DumpZSetValue(v ZSet) []byte {
	buf := bytes.NewBuffer(nil)
	enc := rdb.NewEncoder(buf)
	enc.EncodeType(rdb.TypeZSet)
	enc.EncodeLength(uint32(len(v)))
	for i := 0; i < len(v); i++ {
		enc.EncodeString(v[i].Member)
		enc.EncodeFloat(v[i].Score)
	}
	enc.EncodeDumpFooter()

	return buf.Bytes()
}
