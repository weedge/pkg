package driver

type IWriteBatch interface {
	Put(key []byte, value []byte)
	Delete(key []byte)
	// write to os buffer
	Commit() error
	// sync to stable disk
	SyncCommit() error
	// resets the batch
	Rollback() error
	// batch data
	Data() []byte
	// close WriteBatch
	Close()
}
