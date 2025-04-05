package common

import (
    "sort"
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