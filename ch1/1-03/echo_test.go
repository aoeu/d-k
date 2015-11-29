package ch1

import (
	"fmt"
	"testing"
)

func BenchmarkConcatArgs100(b *testing.B) {
	benchmarkConcatArgs(100, b)
}

func BenchmarkConcatArgs1000(b *testing.B) {
	benchmarkConcatArgs(1000, b)
}

func BenchmarkConcatArgs10000(b *testing.B) {
	benchmarkConcatArgs(10000, b)
}

func benchmarkConcatArgs(n int, b *testing.B) {
	args := make([]string, n)
	for i := 0; i < n; i++ {
		args[i] = fmt.Sprintf("%v", n)
	}
	for i := 0; i < b.N; i++ {
		ConcatArgs(args)
	}
}

func BenchmarkJoinArgs100(b *testing.B) {
	benchmarkJoinArgs(100, b)
}

func BenchmarkJoinArgs1000(b *testing.B) {
	benchmarkJoinArgs(1000, b)
}

func BenchmarkJoinArgs10000(b *testing.B) {
	benchmarkJoinArgs(10000, b)
}

func benchmarkJoinArgs(n int, b *testing.B) {
	args := make([]string, n)
	for i := 0; i < n; i++ {
		args[i] = fmt.Sprintf("%v", n)
	}
	for i := 0; i < b.N; i++ {
		JoinArgs(args)
	}
}
