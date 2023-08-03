package types

import (
	"fmt"
	"io"

	"github.com/weedge/pkg/rdb/structure"
	"github.com/weedge/pkg/utils/logutils"
)

type ZSetEntry struct {
	Member string
	Score  string
}

type ZsetObject struct {
	key      string
	Elements []ZSetEntry
}

func (o *ZsetObject) LoadFromBuffer(rd io.Reader, key string, typeByte byte) {
	switch typeByte {
	case RDBTypeZSet:
		o.readZset(rd)
	case RDBTypeZSet2:
		o.readZset2(rd)
	case RDBTypeZSetZiplist:
		o.readZsetZiplist(rd)
	case RDBTypeZSetListpack:
		o.readZsetListpack(rd)
	default:
		logutils.Criticalf("unknown zset type. typeByte=[%d]", typeByte)
	}
}

func (o *ZsetObject) readZset(rd io.Reader) {
	size := int(structure.ReadLength(rd))
	o.Elements = make([]ZSetEntry, size)
	for i := 0; i < size; i++ {
		o.Elements[i].Member = structure.ReadString(rd)
		score := structure.ReadFloat(rd)
		o.Elements[i].Score = fmt.Sprintf("%f", score)
	}
}

func (o *ZsetObject) readZset2(rd io.Reader) {
	size := int(structure.ReadLength(rd))
	o.Elements = make([]ZSetEntry, size)
	for i := 0; i < size; i++ {
		o.Elements[i].Member = structure.ReadString(rd)
		score := structure.ReadDouble(rd)
		o.Elements[i].Score = fmt.Sprintf("%f", score)
	}
}

func (o *ZsetObject) readZsetZiplist(rd io.Reader) {
	list := structure.ReadZipList(rd)
	size := len(list)
	if size%2 != 0 {
		logutils.Criticalf("zset listpack size is not even. size=[%d]", size)
	}
	o.Elements = make([]ZSetEntry, size/2)
	for i := 0; i < size; i += 2 {
		o.Elements[i/2].Member = list[i]
		o.Elements[i/2].Score = list[i+1]
	}
}

func (o *ZsetObject) readZsetListpack(rd io.Reader) {
	list := structure.ReadListpack(rd)
	size := len(list)
	if size%2 != 0 {
		logutils.Criticalf("zset listpack size is not even. size=[%d]", size)
	}
	o.Elements = make([]ZSetEntry, size/2)
	for i := 0; i < size; i += 2 {
		o.Elements[i/2].Member = list[i]
		o.Elements[i/2].Score = list[i+1]
	}
}

func (o *ZsetObject) Rewrite() []RedisCmd {
	cmds := make([]RedisCmd, len(o.Elements))
	for inx, ele := range o.Elements {
		cmd := RedisCmd{"zadd", o.key, ele.Score, ele.Member}
		cmds[inx] = cmd
	}
	return cmds
}
