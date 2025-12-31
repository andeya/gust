package coarsetime

import (
	"testing"
	"time"
)

// BenchmarkCoarseTime_Now benchmarks the CoarseTime.Now() method
func BenchmarkCoarseTime_Now(b *testing.B) {
	ct := New(100 * time.Millisecond)
	defer ct.Stop()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ct.Now()
	}
}

// BenchmarkTime_Now benchmarks the standard time.Now() for comparison
func BenchmarkTime_Now(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = time.Now()
	}
}

// BenchmarkCoarseTime_Monotonic benchmarks the CoarseTime.Monotonic() method
func BenchmarkCoarseTime_Monotonic(b *testing.B) {
	ct := New(100 * time.Millisecond)
	defer ct.Stop()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ct.Monotonic()
	}
}

// BenchmarkCoarseTime_Since benchmarks the CoarseTime.Since() method
func BenchmarkCoarseTime_Since(b *testing.B) {
	ct := New(100 * time.Millisecond)
	defer ct.Stop()
	start := ct.Monotonic()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ct.Since(start)
	}
}

// BenchmarkTime_Since benchmarks the standard time.Since() for comparison
func BenchmarkTime_Since(b *testing.B) {
	start := time.Now()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = time.Since(start)
	}
}

// BenchmarkConvenience_Now benchmarks the convenience Now() function
func BenchmarkConvenience_Now(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Now()
	}
}

// BenchmarkConvenience_Monotonic benchmarks the convenience Monotonic() function
func BenchmarkConvenience_Monotonic(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Monotonic()
	}
}

// BenchmarkConvenience_Since benchmarks the convenience Since() function
func BenchmarkConvenience_Since(b *testing.B) {
	start := Monotonic()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Since(start)
	}
}

// BenchmarkDifferentFrequencies benchmarks different frequency settings
func BenchmarkDifferentFrequencies(b *testing.B) {
	frequencies := []time.Duration{
		1 * time.Millisecond,
		10 * time.Millisecond,
		100 * time.Millisecond,
		1 * time.Second,
	}

	for _, freq := range frequencies {
		b.Run(freq.String(), func(b *testing.B) {
			ct := New(freq)
			defer ct.Stop()

			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = ct.Now()
			}
		})
	}
}

// BenchmarkConcurrentAccess benchmarks concurrent access to CoarseTime
func BenchmarkConcurrentAccess(b *testing.B) {
	ct := New(100 * time.Millisecond)
	defer ct.Stop()

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = ct.Now()
			_ = ct.Monotonic()
		}
	})
}

// BenchmarkConcurrentTimeNow benchmarks concurrent access to time.Now()
func BenchmarkConcurrentTimeNow(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = time.Now()
		}
	})
}
