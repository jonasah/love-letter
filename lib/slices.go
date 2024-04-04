package lib

func Filter[T any](s []T, f func(T) bool) []T {
	var values []T
	for _, v := range s {
		if f(v) {
			values = append(values, v)
		}
	}

	return values
}

func Transform[T any, U any](s []T, f func(T) U) []U {
	var values []U
	for _, v := range s {
		values = append(values, f(v))
	}

	return values
}
