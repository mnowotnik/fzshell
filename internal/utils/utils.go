package utils

func CloneMap[T comparable, V any](mmap map[T]V) map[T]V {
	clone := make(map[T]V)
	for k, v := range mmap {
		clone[k] = v
	}
	return clone
}
