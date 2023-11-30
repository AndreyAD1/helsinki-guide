package utils

func GetPointer[T string | rune | int | int64 | float32](a T) *T {
	return &a
}
