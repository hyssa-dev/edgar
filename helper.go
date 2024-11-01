package edgar

import "strings"

func hasArray[T comparable](source []T, target T) bool {
	for _, v := range source {
		if v == target {
			return true
		}
	}
	return false
}

func hasNotArray[T comparable](source []T, target T) bool {
	for _, v := range source {
		if v == target {
			return false
		}
	}
	return true
}

func containsAllElements(target string, keys []string) bool {
	cnt := 0
	for _, key := range keys {
		if strings.Contains(target, key) {
			cnt++
		}
	}

	return cnt == len(keys)
}

func containsAnyElement(target string, keys []string) bool {
	for _, key := range keys {
		if strings.Contains(target, key) {
			return true
		}
	}

	return false
}

func nonContainsAllElements(target string, keys []string) bool {
	cnt := 0
	for _, key := range keys {
		if strings.Contains(target, key) {
			cnt++
		}
	}

	return cnt == 0
}
