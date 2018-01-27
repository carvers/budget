package similar

import (
	"fmt"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestShingling(t *testing.T) {
	in := "The cat sat on the cat on the mat"
	expected := []string{"th", "he", "e ", " c", "ca", "at", "t ", " s", "sa", " o", "on", "n ", " t", " m", "ma"}
	result := Shingle(in)
	sort.Strings(expected)
	sort.Strings(result)
	t.Log(expected)
	t.Log(result)
	if !cmp.Equal(expected, result) {
		t.Error(cmp.Diff(expected, result))
	}
}

func TestJaccard(t *testing.T) {
	testpairs := []struct {
		a, b     []string
		expected float64
	}{
		{a: []string{"ab", "bc", "cd"}, b: []string{"ab", "bc", "ce"}, expected: float64(2) / float64(4)},
		{a: []string{"th", "he", "e ", " c", "ca", "at"}, b: []string{"th", "he", "e ", " h", "ha", "at"}, expected: float64(4) / float64(8)},
	}
	for pos, pair := range testpairs {
		t.Run(fmt.Sprintf("pair=%d", pos), func(t *testing.T) {
			output := Jaccard(pair.a, pair.b)
			if output != pair.expected {
				t.Errorf("Expected %f, got %f", pair.expected, output)
			}
		})
	}
}
