// core struct cmd op interface driver
package driver

import (
	"context"
	"time"

	openkvdriver "github.com/weedge/pkg/driver/openkv"
)

type KVPair struct {
	Key   []byte
	Value []byte
}

// adapt https://redis.io/commands/?group=string
type IStringCmd interface {
	Set(ctx context.Context, key []byte, value []byte) error
	SetNX(ctx context.Context, key []byte, value []byte) (n int64, err error)
	SetEX(ctx context.Context, key []byte, duration int64, value []byte) error
	SetNXEX(ctx context.Context, key []byte, duration int64, value []byte) error
	SetXXEX(ctx context.Context, key []byte, duration int64, value []byte) error

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
//	RangeLpen: (min, max)
//	RangeRopen: (min, max]
//	RangeOopen: [min, max)
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

	ICommonCmd
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

type IStorager interface {
	Select(ctx context.Context, index int) (db IDB, err error)
	FlushAll(ctx context.Context) error
	Close() error
}
