package types

import (
	"io"

	"github.com/weedge/pkg/rdb/structure"
	"github.com/weedge/pkg/utils/logutils"
)

type SetObject struct {
	key      string
	Elements []string
}

func (o *SetObject) LoadFromBuffer(rd io.Reader, key string, typeByte byte) {
	o.key = key
	switch typeByte {
	case RDBTypeSet:
		o.readSet(rd)
	case RDBTypeSetIntset:
		o.Elements = structure.ReadIntset(rd)
	default:
		logutils.Criticalf("unknown set type. typeByte=[%d]", typeByte)
	}
}

func (o *SetObject) readSet(rd io.Reader) {
	size := int(structure.ReadLength(rd))
	o.Elements = make([]string, size)
	for i := 0; i < size; i++ {
		val := structure.ReadString(rd)
		o.Elements[i] = val
	}
}

func (o *SetObject) Rewrite() []RedisCmd {
	cmds := make([]RedisCmd, len(o.Elements))
	for inx, ele := range o.Elements {
		cmd := RedisCmd{"sadd", o.key, ele}
		cmds[inx] = cmd
	}
	return cmds
}
