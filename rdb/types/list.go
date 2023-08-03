package types

import (
	"io"

	"github.com/weedge/pkg/rdb/structure"
	"github.com/weedge/pkg/utils/logutils"
)

// quicklist node container formats
const (
	quicklistNodeContainerPlain  = 1 // QUICKLIST_NODE_CONTAINER_PLAIN
	quicklistNodeContainerPacked = 2 // QUICKLIST_NODE_CONTAINER_PACKED
)

type ListObject struct {
	key string

	Elements []string
}

func (o *ListObject) LoadFromBuffer(rd io.Reader, key string, typeByte byte) {
	o.key = key
	switch typeByte {
	case RDBTypeList:
		o.readList(rd)
	case RDBTypeListZiplist:
		o.Elements = structure.ReadZipList(rd)
	case RDBTypeListQuicklist:
		o.readQuickList(rd)
	case RDBTypeListQuicklist2:
		o.readQuickList2(rd)
	default:
		logutils.Criticalf("unknown list type %d", typeByte)
	}
}

func (o *ListObject) Rewrite() []RedisCmd {
	cmds := make([]RedisCmd, len(o.Elements))
	for inx, ele := range o.Elements {
		cmd := RedisCmd{"rpush", o.key, ele}
		cmds[inx] = cmd
	}
	return cmds
}

func (o *ListObject) readList(rd io.Reader) {
	size := int(structure.ReadLength(rd))
	for i := 0; i < size; i++ {
		ele := structure.ReadString(rd)
		o.Elements = append(o.Elements, ele)
	}
}

func (o *ListObject) readQuickList(rd io.Reader) {
	size := int(structure.ReadLength(rd))
	for i := 0; i < size; i++ {
		ziplistElements := structure.ReadZipList(rd)
		o.Elements = append(o.Elements, ziplistElements...)
	}
}

func (o *ListObject) readQuickList2(rd io.Reader) {
	size := int(structure.ReadLength(rd))
	for i := 0; i < size; i++ {
		container := structure.ReadLength(rd)
		if container == quicklistNodeContainerPlain {
			ele := structure.ReadString(rd)
			o.Elements = append(o.Elements, ele)
		} else if container == quicklistNodeContainerPacked {
			listpackElements := structure.ReadListpack(rd)
			o.Elements = append(o.Elements, listpackElements...)
		} else {
			logutils.Criticalf("unknown quicklist container %d", container)
		}
	}
}
