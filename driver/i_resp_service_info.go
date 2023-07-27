package driver

import (
	"fmt"
	"io"
	"strings"

	"github.com/weedge/pkg/utils"
)

type ISrvInfo interface {
	DumpBytes(name DumpSrvInfoName) []byte
}

type DumpSrvInfoName string

func (name DumpSrvInfoName) RespDumpName() []byte {
	nameStr := fmt.Sprintf("# %s\r\n", name.FirstToUp())
	return utils.String2Bytes(nameStr)
}

func (name DumpSrvInfoName) ToLow() []byte {
	return utils.String2Bytes(strings.ToLower(string(name)))
}

func (name DumpSrvInfoName) FirstToUp() []byte {
	str := strings.ToUpper(string(name)[:1]) + string(name)[1:]
	return utils.String2Bytes(str)
}

type InfoPair struct {
	Key   string
	Value interface{}
}

func (pair InfoPair) RespDumpInfo() []byte {
	pairInfo := fmt.Sprintf("%s:%v\r\n", pair.Key, pair.Value)
	return utils.String2Bytes(pairInfo)
}

type DumpHandler func(w io.Writer)

var RegisteredDumpHandlers = map[DumpSrvInfoName]DumpHandler{}
var RegisteredDumpHandlerNames = []DumpSrvInfoName{}

func RegisterDumpHandler(name DumpSrvInfoName, handler DumpHandler) {
	if _, ok := RegisteredDumpHandlers[name]; !ok {
		RegisteredDumpHandlerNames = append(RegisteredDumpHandlerNames, name)
	}
	RegisteredDumpHandlers[name] = handler
}
