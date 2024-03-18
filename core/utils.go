package core

func find[T any](arr []T, predicate func(item T) bool) (T, bool) {
	for _, item := range arr {
		if predicate(item) {
			return item, true
		}
	}
	var result T
	return result, false
}
