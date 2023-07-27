package utils

import "golang.org/x/exp/constraints"

// 利用泛型比较2个数的大小
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T constraints.Ordered](a, b T) T {

	if a < b {
		return a
	}
	return b
}
