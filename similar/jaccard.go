package similar

import "strings"

func Shingle(in string) []string {
	runes := strings.Split(strings.ToLower(in), "")
	shingles := map[string]struct{}{}
	for pos, c := range runes {
		if pos+1 == len(runes) {
			break
		}
		shingles[c+runes[pos+1]] = struct{}{}
	}
	results := make([]string, 0, len(shingles))
	for shingle := range shingles {
		results = append(results, shingle)
	}
	return results
}

func Jaccard(a, b []string) float64 {
	aset := make(map[string]struct{}, len(a))
	union := make(map[string]struct{}, len(a))
	intersection := map[string]struct{}{}
	for _, x := range a {
		aset[x] = struct{}{}
		union[x] = struct{}{}
	}
	for _, x := range b {
		union[x] = struct{}{}
		if _, ok := aset[x]; ok {
			intersection[x] = struct{}{}
		}
	}
	return float64(len(intersection)) / float64(len(union))
}
