package types

import (
	"io"

	"github.com/weedge/pkg/rdb/structure"
	"github.com/weedge/pkg/utils/logutils"
)

type ModuleObject struct {
}

func (o *ModuleObject) LoadFromBuffer(rd io.Reader, key string, typeByte byte) {
	if typeByte == rdbTypeModule {
		logutils.Criticalf("module type with version 1 is not supported, key=[%s]", key)
	}
	moduleId := structure.ReadLength(rd)
	moduleName := moduleTypeNameByID(moduleId)
	opcode := structure.ReadByte(rd)
	for opcode != rdbModuleOpcodeEOF {
		switch opcode {
		case rdbModuleOpcodeSINT:
		case rdbModuleOpcodeUINT:
			structure.ReadLength(rd)
		case rdbModuleOpcodeFLOAT:
			structure.ReadFloat(rd)
		case rdbModuleOpcodeDOUBLE:
			structure.ReadDouble(rd)
		case rdbModuleOpcodeSTRING:
			structure.ReadString(rd)
		default:
			logutils.Criticalf("unknown module opcode=[%d], module name=[%s]", opcode, moduleName)
		}
		opcode = structure.ReadByte(rd)
	}
}

func (o *ModuleObject) Rewrite() []RedisCmd {
	logutils.Criticalf("module Rewrite not implemented")
	return nil
}
