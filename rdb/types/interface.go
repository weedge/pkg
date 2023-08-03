package types

import (
	"io"

	"github.com/weedge/pkg/utils/logutils"
)

const (
	// StringType is redis string
	StringType = "string"
	// ListType is redis list
	ListType = "list"
	// SetType is redis set
	SetType = "set"
	// HashType is redis hash
	HashType = "hash"
	// ZSetType is redis sorted set
	ZSetType = "zset"
	// AuxType is redis metadata key-value pair
	AuxType = "aux"
	// DBSizeType is for _OPCODE_RESIZEDB
	DBSizeType = "dbsize"
)

const (
	RDBTypeString  = 0 // RDB_TYPE_STRING
	RDBTypeList    = 1
	RDBTypeSet     = 2
	RDBTypeZSet    = 3
	RDBTypeHash    = 4 // RDB_TYPE_HASH
	RDBTypeZSet2   = 5 // ZSET version 2 with doubles stored in binary.
	RDBTypeModule  = 6 // RDB_TYPE_MODULE
	RDBTypeModule2 = 7 // RDB_TYPE_MODULE2 Module value with annotations for parsing without the generating module being loaded.

	// Object types for encoded objects.

	RDBTypeHashZipmap       = 9
	RDBTypeListZiplist      = 10
	RDBTypeSetIntset        = 11
	RDBTypeZSetZiplist      = 12
	RDBTypeHashZiplist      = 13
	RDBTypeListQuicklist    = 14 // RDB_TYPE_LIST_QUICKLIST
	RDBTypeStreamListpacks  = 15 // RDB_TYPE_STREAM_LISTPACKS
	RDBTypeHashListpack     = 16 // RDB_TYPE_HASH_ZIPLIST
	RDBTypeZSetListpack     = 17 // RDB_TYPE_ZSET_LISTPACK
	RDBTypeListQuicklist2   = 18 // RDB_TYPE_LIST_QUICKLIST_2 https://github.com/redis/redis/pull/9357
	RDBTypeStreamListpacks2 = 19 // RDB_TYPE_STREAM_LISTPACKS2

	moduleTypeNameCharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

	rdbModuleOpcodeEOF    = 0 // End of module value.
	rdbModuleOpcodeSINT   = 1 // Signed integer.
	rdbModuleOpcodeUINT   = 2 // Unsigned integer.
	rdbModuleOpcodeFLOAT  = 3 // Float.
	rdbModuleOpcodeDOUBLE = 4 // Double.
	rdbModuleOpcodeSTRING = 5 // String.
)

type RedisCmd []string

// RedisObject is interface for a redis object
type RedisObject interface {
	LoadFromBuffer(rd io.Reader, key string, typeByte byte)
	Rewrite() []RedisCmd
}

func ParseObject(rd io.Reader, typeByte byte, key string) RedisObject {
	switch typeByte {
	case RDBTypeString: // string
		o := new(StringObject)
		o.LoadFromBuffer(rd, key, typeByte)
		return o
	case RDBTypeList, RDBTypeListZiplist, RDBTypeListQuicklist, RDBTypeListQuicklist2: // list
		o := new(ListObject)
		o.LoadFromBuffer(rd, key, typeByte)
		return o
	case RDBTypeSet, RDBTypeSetIntset: // set
		o := new(SetObject)
		o.LoadFromBuffer(rd, key, typeByte)
		return o
	case RDBTypeZSet, RDBTypeZSet2, RDBTypeZSetZiplist, RDBTypeZSetListpack: // zset
		o := new(ZsetObject)
		o.LoadFromBuffer(rd, key, typeByte)
		return o
	case RDBTypeHash, RDBTypeHashZipmap, RDBTypeHashZiplist, RDBTypeHashListpack: // hash
		o := new(HashObject)
		o.LoadFromBuffer(rd, key, typeByte)
		return o
	case RDBTypeStreamListpacks, RDBTypeStreamListpacks2: // stream
		o := new(StreamObject)
		o.LoadFromBuffer(rd, key, typeByte)
		return o
	case RDBTypeModule, RDBTypeModule2: // module
		o := new(ModuleObject)
		o.LoadFromBuffer(rd, key, typeByte)
		return o
	}
	logutils.Criticalf("unknown type byte: %d", typeByte)
	return nil
}

func moduleTypeNameByID(moduleId uint64) string {
	nameList := make([]byte, 9)
	moduleId >>= 10
	for i := 8; i >= 0; i-- {
		nameList[i] = moduleTypeNameCharSet[moduleId&63]
		moduleId >>= 6
	}
	return string(nameList)
}
