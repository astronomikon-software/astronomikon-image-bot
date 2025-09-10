package util

func Map[T any, V any](s []T, mapper func(t T) V) []V {
	out := make([]V, 0, len(s))
	for _, item := range s {
		out = append(out, mapper(item))
	}
	return out
}
