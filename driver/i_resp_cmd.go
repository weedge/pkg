// core struct cmd op interface driver
package driver

import (
	"time"

	openkvdriver "github.com/weedge/pkg/driver/openkv"
)

type KVPair struct {
	Key   []byte
	Value []byte
}

// adapt https://redis.io/commands/?group=string
type IStringCmd interface {
	Set(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	GetSlice(key []byte) (openkvdriver.ISlice, error)
	GetSet(key []byte, value []byte) ([]byte, error)
	Incr(key []byte) (int64, error)
	IncrBy(key []byte, increment int64) (int64, error)
	Decr(key []byte) (int64, error)
	DecrBy(key []byte, decrement int64) (int64, error)
	MGet(keys ...[]byte) ([][]byte, error)
	MSet(args ...KVPair) error
	SetNX(key []byte, value []byte) (n int64, err error)
	SetEX(key []byte, duration int64, value []byte) error
	SetRange(key []byte, offset int, value []byte) (int64, error)
	GetRange(key []byte, start int, end int) ([]byte, error)
	StrLen(key []byte) (int64, error)
	Append(key []byte, value []byte) (int64, error)

	ICommonCmd
}

// adapt https://redis.io/commands/?group=list
type IListCmd interface {
	LIndex(key []byte, index int32) ([]byte, error)
	LLen(key []byte) (int64, error)
	LPop(key []byte) ([]byte, error)
	LTrim(key []byte, start, stop int64) error
	LTrimFront(key []byte, trimSize int32) (int32, error)
	LTrimBack(key []byte, trimSize int32) (int32, error)
	LPush(key []byte, args ...[]byte) (int64, error)
	LSet(key []byte, index int32, value []byte) error
	LRange(key []byte, start int32, stop int32) ([][]byte, error)
	RPop(key []byte) ([]byte, error)
	RPush(key []byte, args ...[]byte) (int64, error)
	BLPop(keys [][]byte, timeout time.Duration) ([]interface{}, error)
	BRPop(keys [][]byte, timeout time.Duration) ([]interface{}, error)

	ICommonCmd
}

type FVPair struct {
	Field []byte
	Value []byte
}

// adapt https://redis.io/commands/?group=hash
type IHashCmd interface {
	HSet(key []byte, field []byte, value []byte) (int64, error)
	HGet(key []byte, field []byte) ([]byte, error)
	HLen(key []byte) (int64, error)
	HMset(key []byte, args ...FVPair) error
	HMget(key []byte, args ...[]byte) ([][]byte, error)
	HDel(key []byte, args ...[]byte) (int64, error)
	HIncrBy(key []byte, field []byte, delta int64) (int64, error)
	HGetAll(key []byte) ([]FVPair, error)
	HKeys(key []byte) ([][]byte, error)
	HValues(key []byte) ([][]byte, error)

	ICommonCmd
}

// adapt https://redis.io/commands/?group=set
type ISetCmd interface {
	SAdd(key []byte, args ...[]byte) (int64, error)
	SCard(key []byte) (int64, error)
	SDiff(keys ...[]byte) ([][]byte, error)
	SDiffStore(dstKey []byte, keys ...[]byte) (int64, error)
	SInter(keys ...[]byte) ([][]byte, error)
	SInterStore(dstKey []byte, keys ...[]byte) (int64, error)
	SIsMember(key []byte, member []byte) (int64, error)
	SMembers(key []byte) ([][]byte, error)
	SRem(key []byte, args ...[]byte) (int64, error)
	SUnion(keys ...[]byte) ([][]byte, error)
	SUnionStore(dstKey []byte, keys ...[]byte) (int64, error)

	ICommonCmd
}

type ScorePair struct {
	Score  int64
	Member []byte
}

// adapt https://redis.io/commands/?group=sorted-set
type IZsetCmd interface {
	ZAdd(key []byte, args ...ScorePair) (int64, error)
	ZCard(key []byte) (int64, error)
	ZScore(key []byte, member []byte) (int64, error)
	ZRem(key []byte, members ...[]byte) (int64, error)
	ZIncrBy(key []byte, delta int64, member []byte) (int64, error)
	ZCount(key []byte, min int64, max int64) (int64, error)
	ZRank(key []byte, member []byte) (int64, error)
	ZRemRangeByRank(key []byte, start int, stop int)
	ZRemRangeByScore(key []byte, min int64, max int64) (int64, error)
	ZRevRange(key []byte, start int, stop int) ([]ScorePair, error)
	ZRevRank(key []byte, member []byte) (int64, error)
	ZRevRangeByScore(key []byte, min int64, max int64, offset int, count int)
	ZRangeGeneric(key []byte, start int, stop int, reverse bool) ([]ScorePair, error)
	ZRangeByScoreGeneric(key []byte, min int64, max int64, offset int, count int, reverse bool) ([]ScorePair, error)
	ZUnionStore(destKey []byte, srcKeys [][]byte, weights []int64, aggregate byte) (int64, error)
	ZInterStore(destKey []byte, srcKeys [][]byte, weights []int64, aggregate byte) (int64, error)
	ZRangeByLex(key []byte, min []byte, max []byte, rangeType uint8, offset int, count int) ([][]byte, error)
	ZRemRangeByLex(key []byte, min []byte, max []byte, rangeType uint8) (int64, error)
	ZLexCount(key []byte, min []byte, max []byte, rangeType uint8) (int64, error)

	ICommonCmd
}

// adapt https://redis.io/commands/?group=bitmap
type IBitmapCmd interface {
	BitOP(op string, destKey []byte, srcKeys ...[]byte) (int64, error)
	BitCount(key []byte, start int, end int) (int64, error)
	BitPos(key []byte, on int, start int, end int) (int64, error)
	SetBit(key []byte, offset int, on int) (int64, error)
	GetBit(key []byte, offset int) (int64, error)

	ICommonCmd
}

// adapt https://redis.io/commands/?group=generic
// some common key op cmd
type ICommonCmd interface {
	Del(keys ...[]byte) (int64, error)
	Exists(key []byte) (int64, error)
	// for ttl
	Expire(key []byte, duration int64) (int64, error)
	ExpireAt(key []byte, when int64) (int64, error)
	TTL(key []byte) (int64, error)
	Persist(key []byte) (int64, error)
}

type IExtraCmd interface {
	Clear(key []byte) (int64, error)
	Mclear(keys ...[]byte) (int64, error)
}
