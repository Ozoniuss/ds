package rbtree

import (
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
)

// Run with:
//
//	go test -run '^$' -bench '^Benchmark' -benchmem -count 5
//
// As a comparison, we used an array which is re-sorted after every operation.

var randomValues1_000 = rand.New(rand.NewPCG(6, 9)).Perm(1_000)
var randomValues10_000 = rand.New(rand.NewPCG(6, 9)).Perm(10_000)
var randomValues100_000 = rand.New(rand.NewPCG(6, 9)).Perm(100_000)

var randomValues = map[int][]int{
	1_000:   randomValues1_000,
	10_000:  randomValues10_000,
	100_000: randomValues100_000,
}

func BenchmarkSliceIncreasingSorted(b *testing.B) {
	for n := 1_000; n <= 100_000; n *= 10 {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			for b.Loop() {
				var s []int
				for i := range n {
					s = append(s, i)
					slices.Sort(s)
				}
			}
			b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*n), "ns/insert")
			b.ReportMetric(float64(b.N*n)/b.Elapsed().Seconds(), "inserts/sec")
		})
	}
}

func BenchmarkSliceDecreasingSorted(b *testing.B) {
	for n := 1_000; n <= 100_000; n *= 10 {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			for b.Loop() {
				var s []int
				for i := range n {
					s = append(s, n-1-i)
					slices.Sort(s)
				}
			}
			b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*n), "ns/insert")
			b.ReportMetric(float64(b.N*n)/b.Elapsed().Seconds(), "inserts/sec")
		})
	}
}

func BenchmarkRBTInsertIncreasingSorted(b *testing.B) {
	for n := 1_000; n <= 100_000; n *= 10 {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			for b.Loop() {
				tr := NewRBT[int]()
				for i := range n {
					tr.Insert(i)
				}
			}
			b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*n), "ns/insert")
			b.ReportMetric(float64(b.N*n)/b.Elapsed().Seconds(), "inserts/sec")
		})
	}
}

func BenchmarkSliceSortRandom(b *testing.B) {
	for n := 1_000; n <= 100_000; n *= 10 {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			for b.Loop() {
				var s []int
				for _, v := range randomValues[n] {
					s = append(s, v)
					slices.Sort(s)
				}
			}
			b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*n), "ns/insert")
			b.ReportMetric(float64(b.N*n)/b.Elapsed().Seconds(), "inserts/sec")
		})
	}
}

func BenchmarkRBTInsertRandom(b *testing.B) {
	for n := 1_000; n <= 100_000; n *= 10 {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			for b.Loop() {
				tr := NewRBT[int]()
				for _, v := range randomValues[n] {
					tr.Insert(v)
				}
			}
			b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*n), "ns/insert")
			b.ReportMetric(float64(b.N*n)/b.Elapsed().Seconds(), "inserts/sec")
		})
	}
}
