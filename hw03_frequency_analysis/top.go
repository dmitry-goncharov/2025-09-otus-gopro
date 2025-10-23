package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

const (
	ten   = 10
	cr    = "\n"
	tab   = "\t"
	space = " "
	empty = ""
)

func Top10(str string) []string {
	str = strings.ReplaceAll(str, cr, space)
	str = strings.ReplaceAll(str, tab, space)

	words := strings.Split(str, space)

	m := make(map[string]int, len(words))
	for _, word := range words {
		m[word]++
	}

	delete(m, empty)

	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		if m[keys[i]] == m[keys[j]] {
			return keys[i] < keys[j]
		}
		return m[keys[i]] > m[keys[j]]
	})

	if len(keys) > ten {
		return keys[:ten]
	}

	return keys
}
