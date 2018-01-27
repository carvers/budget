package similar

import (
	"math/rand"
	"testing"
	"time"

	"github.com/tjarratt/babble"
)

func randomWords(babbler babble.Babbler) string {
	babbler.Separator = " "
	rand.Seed(time.Now().UnixNano())
	babbler.Count = rand.Intn(5) + 5
	return babbler.Babble()
}

func init() {
	babbler := babble.NewBabbler()
	// we have roughly 60 transactions/mo
	// five years is 3600 then
	// round up to 4000 and we have a good baseline
	for i := 0; i < 4000; i++ {
		benchDistanceInputs = append(benchDistanceInputs, randomWords(babbler))
	}
}

var benchDistanceInputs []string

func BenchmarkJaccard(b *testing.B) {
	n := len(benchDistanceInputs) - 1
	for i := 0; i < b.N; i++ {
		in1 := benchDistanceInputs[i%n]
		in2 := benchDistanceInputs[(i+1)%n]
		Jaccard(Shingle(in1), Shingle(in2))
	}
}
