package driver

type IIterator interface {
	Close() error

	First()
	Last()
	Seek(key []byte)

	Next()
	Prev()

	Valid() bool

	Key() []byte
	Value() []byte
}
