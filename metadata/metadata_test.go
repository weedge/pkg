package metadata

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPairsMD(t *testing.T) {
	for _, test := range []struct {
		// input
		kv []interface{}
		// output
		md MD
	}{
		{[]interface{}{}, MD{}},
		{[]interface{}{"k1", "v1", "k1", "v2"}, MD{"k1": "v2"}},
	} {
		md := Pairs(test.kv...)
		if !reflect.DeepEqual(md, test.md) {
			t.Fatalf("Pairs(%v) = %v, want %v", test.kv, md, test.md)
		}
	}
}
func TestCopy(t *testing.T) {
	const key, val = "key", "val"
	orig := Pairs(key, val)
	copy := orig.Copy()
	if !reflect.DeepEqual(orig, copy) {
		t.Errorf("copied value not equal to the original, got %v, want %v", copy, orig)
	}
	orig[key] = "foo"
	if v := copy[key]; v != val {
		t.Errorf("change in original should not affect copy, got %q, want %q", v, val)
	}
}
func TestJoin(t *testing.T) {
	for _, test := range []struct {
		mds  []MD
		want MD
	}{
		{[]MD{}, MD{}},
		{[]MD{Pairs("foo", "bar")}, Pairs("foo", "bar")},
		{[]MD{Pairs("foo", "bar"), Pairs("foo", "baz")}, Pairs("foo", "bar", "foo", "baz")},
		{[]MD{Pairs("foo", "bar"), Pairs("foo", "baz"), Pairs("zip", "zap")}, Pairs("foo", "bar", "foo", "baz", "zip", "zap")},
	} {
		md := Join(test.mds...)
		if !reflect.DeepEqual(md, test.want) {
			t.Errorf("context's metadata is %v, want %v", md, test.want)
		}
	}
}

func TestWithContext(t *testing.T) {
	md := MD(map[string]interface{}{"remoteIP": "127.0.0.1", "color": "red", "mirror": true})
	c := NewContext(context.Background(), md)
	ctx := WithContext(c)
	md1, ok := FromContext(ctx)
	if !ok {
		t.Errorf("expect ok be true")
		t.FailNow()
	}
	if !reflect.DeepEqual(md1, md) {
		t.Errorf("expect md1 equal to md")
		t.FailNow()
	}
}

func TestBool(t *testing.T) {
	md := MD{"remoteIP": "127.0.0.1", "color": "red"}
	mdcontext := NewContext(context.Background(), md)
	assert.Equal(t, false, Bool(mdcontext, "mirror"))

	mdcontext = NewContext(context.Background(), MD{"mirror": true})
	assert.Equal(t, true, Bool(mdcontext, "mirror"))

	mdcontext = NewContext(context.Background(), MD{"mirror": "true"})
	assert.Equal(t, true, Bool(mdcontext, "mirror"))

	mdcontext = NewContext(context.Background(), MD{"mirror": "1"})
	assert.Equal(t, true, Bool(mdcontext, "mirror"))

	mdcontext = NewContext(context.Background(), MD{"mirror": "0"})
	assert.Equal(t, false, Bool(mdcontext, "mirror"))

}
func TestInt64(t *testing.T) {
	mdcontext := NewContext(context.Background(), MD{"uid": int64(1)})
	assert.Equal(t, int64(1), Int64(mdcontext, "uid"))
	mdcontext = NewContext(context.Background(), MD{"uid": int64(2)})
	assert.NotEqual(t, int64(1), Int64(mdcontext, "uid"))
	mdcontext = NewContext(context.Background(), MD{"uid": 10})
	assert.NotEqual(t, int64(10), Int64(mdcontext, "uid"))
}
