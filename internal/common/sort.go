package common

import (
	"fmt"
	"sort"
)

func SortKey[T comparable](m map[string]T) []string {
	type kv struct {
		Key   string
		Value T
	}

	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Key < ss[j].Key
	})

	var sortedK []string
	for _, kv := range ss {
		sortedK = append(sortedK, kv.Key)
	}
	return sortedK
}

func SortValue[T any](m map[string]T, compare func(i, j T) bool) []struct {
	Key   string
	Value T
} {
	var ss []struct {
		Key   string
		Value T
	}
	for k, v := range m {
		ss = append(ss, struct {
			Key   string
			Value T
		}{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return compare(ss[i].Value, ss[j].Value)
	})
	return ss
}



func FilterValue[T any](m map[string]T, condition func(T) bool) map[string]T {
	result := make(map[string]T)
	for key, value := range m {
		if condition(value) {
			result[key] = value
		}
	}
	return result
}

func FilterKey(m map[string]int, condition func(string) bool) map[string]int {
	result := make(map[string]int)
	for key, value := range m {
		if condition(key) {
			result[key] = value
		}
	}
	return result
}

// func UnionOperation(a, b []models.Operation) []models.Operation {
// 	seen := make(map[string]bool)
// 	var result []models.Operation

// 	for _, p := range a {
// 		str := fmt.Sprintf("%v", p)
// 		if !seen[str] {
// 			seen[str] = true
// 			result = append(result, p)
// 		}
// 	}

// 	for _, p := range b {
// 		str := fmt.Sprintf("%v", p)
// 		if !seen[str] {
// 			seen[str] = true
// 			result = append(result, p)
// 		}
// 	}

// 	return result
// }


func UnionOperation[T any](a, b []T) []T {
	seen := make(map[string]bool)
	var result []T

	for _, p := range a {
		str := fmt.Sprintf("%v", p)
		if !seen[str] {
			seen[str] = true
			result = append(result, p)
		}
	}

	for _, p := range b {
		str := fmt.Sprintf("%v", p)
		if !seen[str] {
			seen[str] = true
			result = append(result, p)
		}
	}

	return result
}
