// core struct cmd op interface driver
package driver

import (
	"context"
	"fmt"
	"time"

	openkvdriver "github.com/weedge/pkg/driver/openkv"
)

const (
	CmdTypeSrv     = "srv"
	CmdTypeReplica = "replica"
	CmdTypeBitmap  = "bitmap"
	CmdTypeString  = "string"
	CmdTypeHash    = "hash"
	CmdTypeList    = "list"
	CmdTypeSet     = "set"
	CmdTypeZset    = "zset"
	CmdTypeSlot    = "slot"
)

type IReplicaSrvConnCmd interface {
	// Replicaof client cmd in slave
	Replicaof(ctx context.Context, string, restart bool, readonly bool) error
	// Sync internal cmd for slave send sync cmd to master,
	// slave pull buf (sync logs)
	Sync(ctx context.Context, syncLogID uint64) (buf []byte, err error)
	// Sync internal cmd for slave send fullsync cmd to master,
	// slave pull master's snapshot file which dump from data kvstore(FSM),
	// then write to connFD (io.CopyN)
	FullSync(ctx context.Context, needNew bool) (err error)
}

type SlotsRestoreObj struct {
	DB       int32
	Key, Val []byte
	TTLms    int64
}

type SlotInfo struct {
	Num  uint64
	Size uint64
}

type ISlotsCmd interface {
	// MigrateSlotOneKey migrate slot one key/val to addr with timeout (ms)
	// return 1, success, 0 slot is empty
	MigrateSlotOneKey(ctx context.Context, addr string, timeout time.Duration, slot uint64) (int64, error)
	// MigrateSlotKeyWithSameTag migrate slot keys/vals  which have the same tag with one key to addr with timeout (ms)
	// return n, success, 0 slot is empty
	MigrateSlotKeyWithSameTag(ctx context.Context, addr string, timeout time.Duration, slot uint64) (int64, error)
	// MigrateOneKey migrate one key/val (no hash tag  tag=key) to addr with timeout (ms)
	// return n (same key, diff dataType), success, 0 slot is empty
	MigrateOneKey(ctx context.Context, addr string, timeout time.Duration, key []byte) (int64, error)
	// MigrateKeyWithSameTag migrate keys/vals which have the same tag with one key to addr with timeout (ms)
	// return n, n migrate success, 0 slot is empty
	MigrateKeyWithSameTag(ctx context.Context, addr string, timeout time.Duration, key []byte) (int64, error)
	// SlotsRestore dest migrate addr restore slot obj [key ttlms serialized-value(rdb) ...]
	SlotsRestore(ctx context.Context, objs ...*SlotsRestoreObj) error
	// SlotsInfo show slot info with slots range [start,start+count]
	// return slotInfo slice
	SlotsInfo(ctx context.Context, startSlot, count uint64) ([]*SlotInfo, error)
	// SlotsHashKey hash keys to slots, return slot slice
	SlotsHashKey(ctx context.Context, keys ...[]byte) ([]uint64, error)
	// SlotsDel del slots, return after del slot info
	SlotsDel(ctx context.Context, slots ...uint64) ([]*SlotInfo, error)
	// SlotsCheck slots  must check below case
	// - The key stored in each slot can find the corresponding val in the db
	// - Keys in each db can be found in the corresponding slot
	// WARNING: just used debug/test, don't use in product,
	SlotsCheck(ctx context.Context) error
}

type KVPair struct {
	Key   []byte
	Value []byte
}

// adapt https://redis.io/commands/?group=string
type IStringCmd interface {
	Set(ctx context.Context, key []byte, value []byte) error
	SetNX(ctx context.Context, key []byte, value []byte) (n int64, err error)
	SetEX(ctx context.Context, key []byte, duration int64, value []byte) error
	SetNXEX(ctx context.Context, key []byte, duration int64, value []byte) (n int64, err error)
	SetXXEX(ctx context.Context, key []byte, duration int64, value []byte) (n int64, err error)

	Get(ctx context.Context, key []byte) ([]byte, error)
	GetSlice(ctx context.Context, key []byte) (openkvdriver.ISlice, error)
	GetSet(ctx context.Context, key []byte, value []byte) ([]byte, error)

	Incr(ctx context.Context, key []byte) (int64, error)
	IncrBy(ctx context.Context, key []byte, increment int64) (int64, error)
	Decr(ctx context.Context, key []byte) (int64, error)
	DecrBy(ctx context.Context, key []byte, decrement int64) (int64, error)

	MGet(ctx context.Context, keys ...[]byte) ([][]byte, error)
	MSet(ctx context.Context, args ...KVPair) error

	SetRange(ctx context.Context, key []byte, offset int, value []byte) (int64, error)
	GetRange(ctx context.Context, key []byte, start int, end int) ([]byte, error)

	StrLen(ctx context.Context, key []byte) (int64, error)
	Append(ctx context.Context, key []byte, value []byte) (int64, error)

	ICommonCmd
}

// adapt https://redis.io/commands/?group=list
type IListCmd interface {
	LIndex(ctx context.Context, key []byte, index int32) ([]byte, error)
	LLen(ctx context.Context, key []byte) (int64, error)
	LPop(ctx context.Context, key []byte) ([]byte, error)
	LTrim(ctx context.Context, key []byte, start, stop int64) error
	LTrimFront(ctx context.Context, key []byte, trimSize int32) (int32, error)
	LTrimBack(ctx context.Context, key []byte, trimSize int32) (int32, error)
	LPush(ctx context.Context, key []byte, args ...[]byte) (int64, error)
	LSet(ctx context.Context, key []byte, index int32, value []byte) error
	LRange(ctx context.Context, key []byte, start int32, stop int32) ([][]byte, error)
	RPop(ctx context.Context, key []byte) ([]byte, error)
	RPush(ctx context.Context, key []byte, args ...[]byte) (int64, error)
	BLPop(ctx context.Context, keys [][]byte, timeout time.Duration) ([]interface{}, error)
	BRPop(ctx context.Context, keys [][]byte, timeout time.Duration) ([]interface{}, error)

	ICommonCmd
}

type FVPair struct {
	Field []byte
	Value []byte
}

// adapt https://redis.io/commands/?group=hash
type IHashCmd interface {
	HSet(ctx context.Context, key []byte, field []byte, value []byte) (int64, error)
	HGet(ctx context.Context, key []byte, field []byte) ([]byte, error)
	HLen(ctx context.Context, key []byte) (int64, error)
	HMset(ctx context.Context, key []byte, args ...FVPair) error
	HMget(ctx context.Context, key []byte, args ...[]byte) ([][]byte, error)
	HDel(ctx context.Context, key []byte, args ...[]byte) (int64, error)
	HIncrBy(ctx context.Context, key []byte, field []byte, delta int64) (int64, error)
	HGetAll(ctx context.Context, key []byte) ([]FVPair, error)
	HKeys(ctx context.Context, key []byte) ([][]byte, error)
	HValues(ctx context.Context, key []byte) ([][]byte, error)

	ICommonCmd
}

// adapt https://redis.io/commands/?group=set
type ISetCmd interface {
	SAdd(ctx context.Context, key []byte, args ...[]byte) (int64, error)
	SCard(ctx context.Context, key []byte) (int64, error)
	SDiff(ctx context.Context, keys ...[]byte) ([][]byte, error)
	SDiffStore(ctx context.Context, dstKey []byte, keys ...[]byte) (int64, error)
	SInter(ctx context.Context, keys ...[]byte) ([][]byte, error)
	SInterStore(ctx context.Context, dstKey []byte, keys ...[]byte) (int64, error)
	SIsMember(ctx context.Context, key []byte, member []byte) (int64, error)
	SMembers(ctx context.Context, key []byte) ([][]byte, error)
	SRem(ctx context.Context, key []byte, args ...[]byte) (int64, error)
	SUnion(ctx context.Context, keys ...[]byte) ([][]byte, error)
	SUnionStore(ctx context.Context, dstKey []byte, keys ...[]byte) (int64, error)

	ICommonCmd
}

type ScorePair struct {
	Score  int64
	Member []byte
}

// RangeType:
//
//	RangeClose: [min, max]
//	RangeLopen: (min, max]
//	RangeRopen: [min, max)
//	RangeOpen: (min, max)
type RangeType uint8

const (
	RangeClose RangeType = 0x00
	RangeLOpen RangeType = 0x01
	RangeROpen RangeType = 0x10
	RangeOpen  RangeType = 0x11
)

// adapt https://redis.io/commands/?group=sorted-set
type IZsetCmd interface {
	ZAdd(ctx context.Context, key []byte, args ...ScorePair) (int64, error)
	ZCard(ctx context.Context, key []byte) (int64, error)
	ZScore(ctx context.Context, key []byte, member []byte) (int64, error)
	ZRem(ctx context.Context, key []byte, members ...[]byte) (int64, error)
	ZIncrBy(ctx context.Context, key []byte, delta int64, member []byte) (int64, error)
	ZCount(ctx context.Context, key []byte, min int64, max int64) (int64, error)
	ZRank(ctx context.Context, key []byte, member []byte) (int64, error)
	ZRemRangeByRank(ctx context.Context, key []byte, start int, stop int) (int64, error)
	ZRemRangeByScore(ctx context.Context, key []byte, min int64, max int64) (int64, error)
	ZRevRange(ctx context.Context, key []byte, start int, stop int) ([]ScorePair, error)
	ZRevRank(ctx context.Context, key []byte, member []byte) (int64, error)
	ZRevRangeByScore(ctx context.Context, key []byte, min int64, max int64, offset int, count int) ([]ScorePair, error)
	ZRangeGeneric(ctx context.Context, key []byte, start int, stop int, reverse bool) ([]ScorePair, error)
	ZRangeByScoreGeneric(ctx context.Context, key []byte, min int64, max int64, offset int, count int, reverse bool) ([]ScorePair, error)
	ZUnionStore(ctx context.Context, destKey []byte, srcKeys [][]byte, weights []int64, aggregate []byte) (int64, error)
	ZInterStore(ctx context.Context, destKey []byte, srcKeys [][]byte, weights []int64, aggregate []byte) (int64, error)
	ZRangeByLex(ctx context.Context, key []byte, min []byte, max []byte, rangeType RangeType, offset int, count int) ([][]byte, error)
	ZRemRangeByLex(ctx context.Context, key []byte, min []byte, max []byte, rangeType RangeType) (int64, error)
	ZLexCount(ctx context.Context, key []byte, min []byte, max []byte, rangeType RangeType) (int64, error)

	ICommonCmd
}

// adapt https://redis.io/commands/?group=bitmap
type IBitmapCmd interface {
	BitOP(ctx context.Context, op string, destKey []byte, srcKeys ...[]byte) (int64, error)
	BitCount(ctx context.Context, key []byte, start int, end int) (int64, error)
	BitPos(ctx context.Context, key []byte, on int, start int, end int) (int64, error)
	SetBit(ctx context.Context, key []byte, offset int, on int) (int64, error)
	GetBit(ctx context.Context, key []byte, offset int) (int64, error)

	//ICommonCmd
}

// adapt https://redis.io/commands/?group=generic
// some common key op cmd
type ICommonCmd interface {
	Del(ctx context.Context, keys ...[]byte) (int64, error)
	Exists(ctx context.Context, key []byte) (int64, error)
	// for ttl
	Expire(ctx context.Context, key []byte, duration int64) (int64, error)
	ExpireAt(ctx context.Context, key []byte, when int64) (int64, error)
	TTL(ctx context.Context, key []byte) (int64, error)
	Persist(ctx context.Context, key []byte) (int64, error)
}

type IDB interface {
	FlushDB(ctx context.Context) (drop int64, err error)

	DBString() IStringCmd
	DBList() IListCmd
	DBHash() IHashCmd
	DBSet() ISetCmd
	DBZSet() IZsetCmd
	DBBitmap() IBitmapCmd
}

type IDBSlots interface {
	IDB
	DBSlot() ISlotsCmd
}

type IStorager interface {
	Select(ctx context.Context, index int) (db IDB, err error)
	FlushAll(ctx context.Context) error
	Open(ctx context.Context) error
	Close() error
	Name() string
}

type IStatsStorager interface {
	IStorager
	StatsInfo(sections ...string) (info map[string][]InfoPair)
}

var storagers = map[string]IStorager{}

func RegisterStorager(s IStorager) error {
	name := s.Name()
	if _, ok := storagers[name]; ok {
		return fmt.Errorf("storager %s is registered", s)
	}

	storagers[name] = s
	return nil
}

func ListStoragers() []string {
	s := []string{}
	for k := range storagers {
		s = append(s, k)
	}

	return s
}

func GetStorager(name string) (IStorager, error) {
	s, ok := storagers[name]
	if !ok {
		return nil, fmt.Errorf("kv storager %s is not registered", name)
	}

	return s, nil
}
