package types

import (
	"io"

	"github.com/weedge/pkg/rdb/structure"
)

type StringObject struct {
	Key   string
	Value string
}

func (o *StringObject) LoadFromBuffer(rd io.Reader, key string, _ byte) {
	o.Key = key
	o.Value = structure.ReadString(rd)
}

func (o *StringObject) Rewrite() []RedisCmd {
	cmd := RedisCmd{}
	cmd = append(cmd, "set", o.Key, o.Value)
	return []RedisCmd{cmd}
}
