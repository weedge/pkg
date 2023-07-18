package driver

import (
	"context"
	"errors"
	"strings"
)

type IRespConn interface {
	SetDb(db IDB)
	Db() (db IDB)
	SetConnName(name string)
	Name() (name string)
	DoCmd(ctx context.Context, cmd string, cmdParams [][]byte) (res interface{}, err error)
	Close() error
}

type CmdHandle func(ctx context.Context, c IRespConn, cmdParams [][]byte) (interface{}, error)

var RegisteredCmdHandles = map[string]CmdHandle{}
var RegisteredReplicaCmdHandles = map[string]CmdHandle{}
var RegisteredCmdSet = map[string][]string{}

// RegisterCmd register all cmd
func RegisterCmd(cmdType, cmd string, handle CmdHandle) {
	if _, ok := RegisteredCmdHandles[cmd]; ok {
		return
	}

	switch cmdType {
	case CmdTypeReplica:
		RegisteredReplicaCmdHandles[cmd] = handle
	default:
		RegisteredCmdHandles[cmd] = handle
	}
	RegisteredCmdSet[cmdType] = append(RegisteredCmdSet[cmdType], cmd)
}

func MergeRegisteredCmdHandles(src, dst map[string]CmdHandle, isDelSrc bool) {
	for k, v := range src {
		if _, ok := dst[k]; !ok {
			dst[k] = v
		}
		if isDelSrc {
			delete(src, k)
		}
	}
}

type RespConnBase struct {
	db   IDB
	name string
}

func (c *RespConnBase) SetDb(db IDB) {
	c.db = db
}
func (c *RespConnBase) Db() (db IDB) {
	return c.db
}

func (c *RespConnBase) SetConnName(name string) {
	c.name = name
}
func (c *RespConnBase) Name() (name string) {
	return c.name
}

func (c *RespConnBase) Close() error {
	return nil
}

func (c *RespConnBase) DoCmd(ctx context.Context, cmd string, cmdParams [][]byte) (res interface{}, err error) {
	cmd = strings.ToLower(strings.TrimSpace(cmd))
	f, ok := RegisteredCmdHandles[cmd]
	if !ok {
		err = errors.New("ERR unknown command '" + cmd + "'")
		return
	}

	res, err = f(ctx, c, cmdParams)
	if err != nil {
		return
	}

	return
}

var RegisteredWriteCmdAtProposeHandles = map[string]CmdHandle{}
var RegisteredReadCmdAtProposeHandles = map[string]CmdHandle{}
var RegisteredWriteCmdAtApplyHandles = map[string]CmdHandle{}
var RegisteredReadCmdAtApplyHandles = map[string]CmdHandle{}

// RegisterWriteCmdAtPropose
func RegisterWriteCmdAtPropose(cmdType, cmd string, handle CmdHandle) {
	if _, ok := RegisteredWriteCmdAtProposeHandles[cmd]; ok {
		return
	}
	RegisteredWriteCmdAtProposeHandles[cmd] = handle
}

// RegisterReadCmdAtPropose
func RegisteredReadCmdAtPropose(cmdType, cmd string, handle CmdHandle) {
	if _, ok := RegisteredReadCmdAtProposeHandles[cmd]; ok {
		return
	}
	RegisteredReadCmdAtProposeHandles[cmd] = handle
}

// RegisterWriteCmdAtApply
func RegisterWriteCmdAtApply(cmdType, cmd string, handle CmdHandle) {
	if _, ok := RegisteredWriteCmdAtApplyHandles[cmd]; ok {
		return
	}
	RegisteredWriteCmdAtApplyHandles[cmd] = handle
}

// RegisterReadCmdAtApply
func RegisteredReadCmdAtApply(cmdType, cmd string, handle CmdHandle) {
	if _, ok := RegisteredReadCmdAtApplyHandles[cmd]; ok {
		return
	}
	RegisteredReadCmdAtApplyHandles[cmd] = handle
}
