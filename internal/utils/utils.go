package utils

func GetPointer[T string | int | int64 | float32](a T) *T {
	return &a
}
