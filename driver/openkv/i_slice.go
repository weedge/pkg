package driver

// ISlice interface for use go/cgo leveldb/rocksdb lib slice op
type ISlice interface {
	Data() []byte
	Size() int
	Free()
}

type GoSlice []byte

func (s GoSlice) Data() []byte {
	return []byte(s)
}

func (s GoSlice) Size() int {
	return len(s)
}

func (s GoSlice) Free() {

}

// ISliceGetrer interface for use cgo leveldb/rocksdb lib slice get op
type ISliceGeter interface {
	GetSlice(key []byte) (ISlice, error)
}
