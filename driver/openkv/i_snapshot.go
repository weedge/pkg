package driver

type ISnapshot interface {
	Get(key []byte) ([]byte, error)
	NewIterator() IIterator
	Close()
}
