package driver

type IWriteBatch interface {
	Put(key []byte, value []byte)
	Delete(key []byte)
	Commit() error
	SyncCommit() error
	Rollback() error
	Data() []byte
	Close()
}
