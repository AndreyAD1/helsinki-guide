package utils

func GetPointer[T string | rune | int | int64 | float32 | float64](a T) *T {
	return &a
}
