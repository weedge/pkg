package types

import (
	"io"

	"github.com/weedge/pkg/rdb/structure"
	"github.com/weedge/pkg/utils/logutils"
)

type HashObject struct {
	key   string
	Value map[string]string
}

func (o *HashObject) LoadFromBuffer(rd io.Reader, key string, typeByte byte) {
	o.key = key
	o.Value = make(map[string]string)
	switch typeByte {
	case RDBTypeHash:
		o.readHash(rd)
	case RDBTypeHashZipmap:
		o.readHashZipmap(rd)
	case RDBTypeHashZiplist:
		o.readHashZiplist(rd)
	case RDBTypeHashListpack:
		o.readHashListpack(rd)
	default:
		logutils.Criticalf("unknown hash type. typeByte=[%d]", typeByte)
	}
}

func (o *HashObject) readHash(rd io.Reader) {
	size := int(structure.ReadLength(rd))
	for i := 0; i < size; i++ {
		key := structure.ReadString(rd)
		value := structure.ReadString(rd)
		o.Value[key] = value
	}
}

func (o *HashObject) readHashZipmap(rd io.Reader) {
	logutils.Criticalf("not implemented RDBTypeZipmap")
}

func (o *HashObject) readHashZiplist(rd io.Reader) {
	list := structure.ReadZipList(rd)
	size := len(list)
	for i := 0; i < size; i += 2 {
		key := list[i]
		value := list[i+1]
		o.Value[key] = value
	}
}

func (o *HashObject) readHashListpack(rd io.Reader) {
	list := structure.ReadListpack(rd)
	size := len(list)
	for i := 0; i < size; i += 2 {
		key := list[i]
		value := list[i+1]
		o.Value[key] = value
	}
}

func (o *HashObject) Rewrite() []RedisCmd {
	var cmds []RedisCmd
	for k, v := range o.Value {
		cmd := RedisCmd{"hset", o.key, k, v}
		cmds = append(cmds, cmd)
	}
	return cmds
}
