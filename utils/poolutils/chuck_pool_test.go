package poolutils

import (
	"testing"
)

// if make []byte <= 64K don't to return, does not escape, just in stack, no GC
// but use chuck pool can cache reuse, reduce op

func TestSyncPoolAllocBadTiny(t *testing.T) {
	//pool := NewSyncPool(18, 8, 2)
	//pool := NewSyncPool(18, 8, 1)
	pool := NewSyncPool(8, 8, 1)
	mem := pool.Alloc(1)
	if len(mem) != 1 {
		t.Fatal("no equal")
	}
	if cap(mem) != 1 {
		t.Fatal("no equal")
	}
	pool.Free(mem)
}

func TestSyncPoolAllocTiny(t *testing.T) {
	pool := NewSyncPool(8, 8, 2)
	mem := pool.Alloc(1)
	if len(mem) != 1 {
		t.Fatal("no equal")
	}
	if cap(mem) != 8 {
		t.Fatal("no equal")
	}
	pool.Free(mem)
}

func TestSyncPoolAllocSmall(t *testing.T) {
	pool := NewSyncPool(128, 1024, 2)
	mem := pool.Alloc(64)
	if len(mem) != 64 {
		t.Fatal("no equal")
	}
	if cap(mem) != 128 {
		t.Fatal("no equal")
	}
	pool.Free(mem)
}

func TestSyncPoolAllocLarge(t *testing.T) {
	pool := NewSyncPool(128, 1024, 2)
	mem := pool.Alloc(2048)
	if len(mem) != 2048 {
		t.Fatal("no equal")
	}
	if cap(mem) != 2048 {
		t.Fatal("no equal")
	}
	pool.Free(mem)
}

func BenchmarkSyncPoolAllocAndFree128(b *testing.B) {
	pool := NewSyncPool(128, 1024, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Free(pool.Alloc(128))
		}
	})
}

func BenchmarkSyncPoolAllocAndFree256(b *testing.B) {
	pool := NewSyncPool(128, 1024, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Free(pool.Alloc(256))
		}
	})
}

func BenchmarkSyncPoolAllocAndFree512(b *testing.B) {
	pool := NewSyncPool(128, 1024, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Free(pool.Alloc(512))
		}
	})
}

func BenchmarkSyncPoolCacheMiss128(b *testing.B) {
	pool := NewSyncPool(128, 1024, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Alloc(128)
		}
	})
}

func Benchmark_SyncPool_CacheMiss_256(b *testing.B) {
	pool := NewSyncPool(128, 1024, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Alloc(256)
		}
	})
}

func Benchmark_SyncPool_CacheMiss_512(b *testing.B) {
	pool := NewSyncPool(128, 1024, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Alloc(512)
		}
	})
}

func Benchmark_Make_128(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var x []byte
		for pb.Next() {
			x = make([]byte, 128)
		}
		_ = x[:0]
	})
}

func Benchmark_Make_256(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var x []byte
		for pb.Next() {
			x = make([]byte, 256)
		}
		_ = x[:0]
	})
}

func Benchmark_Make_512(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var x []byte
		for pb.Next() {
			x = make([]byte, 512)
		}
		_ = x[:0]
	})
}
