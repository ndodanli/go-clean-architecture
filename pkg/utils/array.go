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
