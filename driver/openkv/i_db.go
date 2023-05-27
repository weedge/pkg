package driver

type IDB interface {
	Close() error

	Get(key []byte) ([]byte, error)

	Put(key []byte, value []byte) error
	Delete(key []byte) error

	SyncPut(key []byte, value []byte) error
	SyncDelete(key []byte) error

	NewIterator() IIterator

	NewWriteBatch() IWriteBatch

	NewSnapshot() (ISnapshot, error)

	Compact() error
}
