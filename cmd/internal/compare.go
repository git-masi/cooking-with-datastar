package internal

func Ternary[T any](condition bool, a T, b T) T {
	if condition {
		return a
	} else {
		return b
	}
}
