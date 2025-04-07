package common

import (
	"fmt"
	"log"
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
	seen := make(map[string]bool)            // Using ID as the unique key
	var result []models.Operation

	// Add elements from first slice
	for _, p := range a {
        str := fmt.Sprintf("%v", p)
		if !seen[str] {
			seen[str] = true
			result = append(result, p)
		} else {
			log.Printf("skip 1 %v", p)
		}
	}

	// Add elements from second slice
	for _, p := range b {
        str := fmt.Sprintf("%v", p)
		if !seen[str] {
			seen[str] = true
			result = append(result, p)
		} else {
			log.Printf("skip 2 \n%#v=?=\n", p)
		}
	}

	return result
}
