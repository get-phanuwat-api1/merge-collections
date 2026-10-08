package merge_test

import (
	"testing"

	"github.com/get-phanuwat-api1/merge-collections"
)

// buildAscending turns arbitrary bytes into a valid ascending collection starting from >= 0
// by accumulating non-negative steps (no sorting needed).
func buildAscending(steps []byte) []int {
	out := make([]int, len(steps))
	cur := 0
	for i, s := range steps {
		cur += int(s % 8) // small steps produce many duplicates
		out[i] = cur
	}
	return out
}

// buildDescending builds a valid descending collection (may go negative).
func buildDescending(start int8, steps []byte) []int {
	out := make([]int, len(steps))
	cur := int(start)
	for i, s := range steps {
		out[i] = cur
		cur -= int(s % 8)
	}
	return out
}

// FuzzMerge checks for random valid inputs that the result is ascending and
// contains exactly the input elements (same multiset), without any sorting.
//
//	go test -run '^$' -fuzz FuzzMerge -fuzztime 30s
func FuzzMerge(f *testing.F) {
	f.Add(int8(9), []byte{2, 2, 2}, []byte{0, 2, 2}, []byte{1, 5})
	f.Add(int8(0), []byte{}, []byte{}, []byte{})
	f.Add(int8(-5), []byte{0, 0, 7}, []byte{0, 0}, []byte{})

	f.Fuzz(func(t *testing.T, start int8, s1, s2, s3 []byte) {
		c1 := buildDescending(start, s1)
		c2 := buildAscending(s2)
		c3 := buildAscending(s3)
		if err := merge.Validate(c1, c2, c3); err != nil {
			t.Fatalf("generator produced invalid input: %v", err)
		}

		got := merge.Merge(c1, c2, c3)

		if len(got) != len(c1)+len(c2)+len(c3) {
			t.Fatalf("len = %d, want %d", len(got), len(c1)+len(c2)+len(c3))
		}
		for i := 1; i < len(got); i++ {
			if got[i] < got[i-1] {
				t.Fatalf("not ascending at %d: %v", i, got)
			}
		}
		counts := map[int]int{}
		for _, c := range [][]int{c1, c2, c3} {
			for _, v := range c {
				counts[v]++
			}
		}
		for _, v := range got {
			counts[v]--
		}
		for v, n := range counts {
			if n != 0 {
				t.Fatalf("value %d appears %d times too many/few in %v", v, n, got)
			}
		}
	})
}
