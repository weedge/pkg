package rdb

// Copyright 2014 Wandoujia Inc. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/cupcake/rdb"
	"github.com/cupcake/rdb/nopdecoder"
	"github.com/weedge/pkg/rdb/structure"
	"github.com/weedge/pkg/rdb/types"
	"github.com/weedge/pkg/utils"
)

func DecodeDump(p []byte) (interface{}, error) {
	rd := bytes.NewReader(p)
	typeByte := structure.ReadByte(rd)
	//key := structure.ReadString(rd)
	o := types.ParseObject(rd, typeByte, "")
	switch item := o.(type) {
	case *types.StringObject:
		return String(item.Value), nil
	case *types.HashObject:
		data := make(Hash, len(item.Value))
		i := 0
		for f, v := range item.Value {
			data[i].Field = utils.String2Bytes(f)
			data[i].Value = utils.String2Bytes(v)
			i++
		}
		return data, nil
	case *types.ListObject:
		data := make(List, len(item.Elements))
		for i := range item.Elements {
			data[i] = utils.String2Bytes(item.Elements[i])
			data[i] = utils.String2Bytes(item.Elements[i])
		}
		return data, nil
	case *types.SetObject:
		data := make(Set, len(item.Elements))
		for i := range item.Elements {
			data[i] = utils.String2Bytes(item.Elements[i])
			data[i] = utils.String2Bytes(item.Elements[i])
		}
		return data, nil
	case *types.ZsetObject:
		data := make(ZSet, len(item.Elements))
		for i := range item.Elements {
			data[i].Member = utils.String2Bytes(item.Elements[i].Member)
			if s, err := strconv.ParseFloat(item.Elements[i].Score, 64); err == nil {
				return nil, err
			} else {
				data[i].Score = s
			}
		}
		return data, nil
	}

	return nil, nil
}

func decodeDump(p []byte) (interface{}, error) {
	d := &decoder{}
	if err := rdb.DecodeDump(p, 0, nil, 0, d); err != nil {
		return nil, err
	}
	return d.obj, d.err
}

type decoder struct {
	nopdecoder.NopDecoder
	obj interface{}
	err error
}

func (d *decoder) initObject(obj interface{}) {
	if d.err != nil {
		return
	}
	if d.obj != nil {
		d.err = fmt.Errorf("invalid object, init again")
	} else {
		d.obj = obj
	}
}

func (d *decoder) Set(key, value []byte, expiry int64) {
	d.initObject(String(value))
}

func (d *decoder) StartHash(key []byte, length, expiry int64) {
	d.initObject(Hash(nil))
}

func (d *decoder) Hset(key, field, value []byte) {
	if d.err != nil {
		return
	}
	switch h := d.obj.(type) {
	default:
		d.err = fmt.Errorf("invalid object, not a hashmap")
	case Hash:
		v := struct {
			Field, Value []byte
		}{
			field,
			value,
		}
		d.obj = append(h, v)
	}
}

func (d *decoder) StartSet(key []byte, cardinality, expiry int64) {
	d.initObject(Set(nil))
}

func (d *decoder) Sadd(key, member []byte) {
	if d.err != nil {
		return
	}
	switch s := d.obj.(type) {
	default:
		d.err = fmt.Errorf("invalid object, not a set")
	case Set:
		d.obj = append(s, member)
	}
}

func (d *decoder) StartList(key []byte, length, expiry int64) {
	d.initObject(List(nil))
}

func (d *decoder) Rpush(key, value []byte) {
	if d.err != nil {
		return
	}
	switch l := d.obj.(type) {
	default:
		d.err = fmt.Errorf("invalid object, not a list")
	case List:
		d.obj = append(l, value)
	}
}

func (d *decoder) StartZSet(key []byte, cardinality, expiry int64) {
	d.initObject(ZSet(nil))
}

func (d *decoder) Zadd(key []byte, score float64, member []byte) {
	if d.err != nil {
		return
	}
	switch z := d.obj.(type) {
	default:
		d.err = fmt.Errorf("invalid object, not a zset")
	case ZSet:
		v := struct {
			Member []byte
			Score  float64
		}{
			member,
			score,
		}
		d.obj = append(z, v)
	}
}
