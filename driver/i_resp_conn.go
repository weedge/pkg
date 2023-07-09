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
}

type cmdHandle func(ctx context.Context, h IRespConn, cmdParams [][]byte) (interface{}, error)

var RegisteredCmdHandles = map[string]cmdHandle{}
var RegisteredCmdSet = map[string][]string{}

// RegisterCmd  register all cmd
func RegisterCmd(cmdType, cmd string, handle cmdHandle) {
	if _, ok := RegisteredCmdHandles[cmd]; ok {
		return
	}
	RegisteredCmdHandles[cmd] = handle
	RegisteredCmdSet[cmdType] = append(RegisteredCmdSet[cmdType], cmd)
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

var RegisteredWriteCmdAtProposeHandles = map[string]cmdHandle{}
var RegisteredReadCmdAtProposeHandles = map[string]cmdHandle{}
var RegisteredWriteCmdAtApplyHandles = map[string]cmdHandle{}
var RegisteredReadCmdAtApplyHandles = map[string]cmdHandle{}

// RegisterWriteCmdAtPropose
func RegisterWriteCmdAtPropose(cmdType, cmd string, handle cmdHandle) {
	if _, ok := RegisteredWriteCmdAtProposeHandles[cmd]; ok {
		return
	}
	RegisteredWriteCmdAtProposeHandles[cmd] = handle
}

// RegisterWriteCmdAtPropose
func RegisteredReadCmdAtPropose(cmdType, cmd string, handle cmdHandle) {
	if _, ok := RegisteredReadCmdAtProposeHandles[cmd]; ok {
		return
	}
	RegisteredReadCmdAtProposeHandles[cmd] = handle
}

// RegisterWriteCmdAtApply
func RegisterWriteCmdAtApply(cmdType, cmd string, handle cmdHandle) {
	if _, ok := RegisteredWriteCmdAtApplyHandles[cmd]; ok {
		return
	}
	RegisteredWriteCmdAtApplyHandles[cmd] = handle
}

// RegisterWriteCmdAtApply
func RegisteredReadCmdAtApply(cmdType, cmd string, handle cmdHandle) {
	if _, ok := RegisteredReadCmdAtApplyHandles[cmd]; ok {
		return
	}
	RegisteredReadCmdAtApplyHandles[cmd] = handle
}
