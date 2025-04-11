package common

import (
	"fmt"
	"sort"
	"test2/internal/models"
)

func Sort(m map[string]int) []string {
	type kv struct {
		Key   string
		Value int
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

func UnionOperation(a, b []models.Operation) []models.Operation {
	seen := make(map[string]bool)
	var result []models.Operation

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
