package util

import "golang.org/x/exp/constraints"

func Min[T constraints.Integer](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T constraints.Integer](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Mod[T constraints.Integer](a, b T) T {
	return (a + b) % b
}
