package utils

func ArrayIncludes[T comparable](slice []T, target T) bool {
	for _, value := range slice {
		if value == target {
			return true
		}
	}
	return false
}

func ArrayAll[T any](slice []T, predicate func(T) bool) bool {
	for _, value := range slice {
		if !predicate(value) {
			return false
		}
	}
	return true
}

func ArrayAny[T any](slice []T, predicate func(T) bool) bool {
	for _, value := range slice {
		if predicate(value) {
			return true
		}
	}
	return false
}

func ArrayFilter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, value := range slice {
		if predicate(value) {
			result = append(result, value)
		}
	}
	return result
}

func ArrayMap[T any, U any](slice []T, mapper func(T) U) []U {
	result := make([]U, len(slice))
	for i, value := range slice {
		result[i] = mapper(value)
	}
	return result
}

func ArrayReduce[T any, U any](slice []T, reducer func(U, T) U, initialValue U) U {
	result := initialValue
	for _, value := range slice {
		result = reducer(result, value)
	}
	return result
}

func ArrayFind[T any](slice []T, predicate func(T) bool) *T {
	for _, value := range slice {
		if predicate(value) {
			return &value
		}
	}
	return nil
}

func ArrayFindIndex[T any](slice []T, predicate func(T) bool) int {
	for i, value := range slice {
		if predicate(value) {
			return i
		}
	}
	return -1
}

func ArrayFindLastIndex[T any](slice []T, predicate func(T) bool) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if predicate(slice[i]) {
			return i
		}
	}
	return -1
}

func ArrayForEach[T any](slice []T, callback func(T)) {
	for _, value := range slice {
		callback(value)
	}
}
