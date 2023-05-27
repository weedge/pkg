package driver

import (
	"fmt"
)

type IStore interface {
	String() string
	Open(path string) (IDB, error)
	Repair(path string) error
}

var dbs = map[string]IStore{}

func Register(s IStore) {
	name := s.String()
	if _, ok := dbs[name]; ok {
		panic(fmt.Errorf("store %s is registered", s))
	}

	dbs[name] = s
}

func ListStores() []string {
	s := []string{}
	for k := range dbs {
		s = append(s, k)
	}

	return s
}

func GetStore(kvStoreName string) (IStore, error) {
	s, ok := dbs[kvStoreName]
	if !ok {
		return nil, fmt.Errorf("kv store engine %s is not registered", kvStoreName)
	}

	return s, nil
}
