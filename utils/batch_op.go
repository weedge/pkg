package utils

import "container/list"

type BatchOpBuffer struct {
	// use list to batch once commit batch put one i/o syscall
	// async --> system buff - flush ->disk
	// sync - flush -> disk
	OpList list.List
}

const (
	BatchOpTypeUnkonw byte = iota
	BatchOpTypePut
	BatchOpTypeDel
)

type BatchOp struct {
	Type       byte
	Key, Value []byte
}

func NewBatchOpBuffer() *BatchOpBuffer {
	return &BatchOpBuffer{}
}

func (bt *BatchOpBuffer) Put(key, value []byte) {
	bt.OpList.PushBack(&BatchOp{Key: key, Value: value, Type: BatchOpTypePut})
}

func (bt *BatchOpBuffer) Del(key []byte) {
	bt.OpList.PushBack(&BatchOp{Key: key, Type: BatchOpTypeDel})
}

func (bt *BatchOpBuffer) FrontElement() *list.Element {
	return bt.OpList.Front()
}

func (bt *BatchOpBuffer) Reset() {
	bt.OpList.Init()
}

func (bt *BatchOpBuffer) Len() int {
	return bt.OpList.Len()
}
