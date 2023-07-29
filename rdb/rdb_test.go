package rdb

import (
	"reflect"
	"testing"
)

func TestCodec(t *testing.T) {
	testCodec(String("abc"), t)
}

func testCodec(obj interface{}, t *testing.T) {
	var rdbData []byte
	switch val := obj.(type) {
	case String:
		rdbData = DumpStringValue(val)
	}

	if o, err := DecodeDump(rdbData); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(obj, o) {
		t.Fatal("must equal")
	}
}
