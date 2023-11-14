package main

import "testing"

// ベンチマークテスト
func Benchmark_main(b *testing.B) {
	for i := 0; i < b.N; i++ {
		main()
	}
}

// sonnetパッケージを使った場合のベンチマークテスト
func Benchmark_sonnet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sonnet()
	}
}

// 標準パッケージを使った場合のベンチマークテスト
func Benchmark_std(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Std()
	}
}
