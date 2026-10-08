// Package merge merges three pre-sorted integer collections into one ascending slice
// without using any sort function.
package merge

import (
	"errors"
	"fmt"
)

// Merge returns every element of the three collections in ascending order.
//
// Preconditions (as given by the specification):
//   - collection1 is sorted from max to min (descending)
//   - collection2 and collection3 are sorted from min (0) to max (ascending)
//
// It is a three-way merge: collection1 is read from its end (its smallest element)
// and the other two from their start, always taking the smallest head.
// Time O(n), extra space O(n) for the result only. Inputs are not modified.
func Merge(collection1, collection2, collection3 []int) []int {
	result := make([]int, 0, len(collection1)+len(collection2)+len(collection3))

	i := len(collection1) - 1 // collection1 is descending: walk backwards
	j, k := 0, 0              // collection2 / collection3 are ascending: walk forwards

	for i >= 0 || j < len(collection2) || k < len(collection3) {
		const (
			fromC1 = iota + 1
			fromC2
			fromC3
		)
		source, smallest := 0, 0

		if i >= 0 {
			source, smallest = fromC1, collection1[i]
		}
		if j < len(collection2) && (source == 0 || collection2[j] < smallest) {
			source, smallest = fromC2, collection2[j]
		}
		if k < len(collection3) && (source == 0 || collection3[k] < smallest) {
			source, smallest = fromC3, collection3[k]
		}

		result = append(result, smallest)
		switch source {
		case fromC1:
			i--
		case fromC2:
			j++
		case fromC3:
			k++
		}
	}
	return result
}

// ErrInvalidInput is wrapped by Validate when a collection breaks the preconditions.
var ErrInvalidInput = errors.New("invalid input")

// Validate checks the preconditions of Merge. Merge itself trusts its input (to stay O(n)
// with no extra checks); call Validate first when the input comes from outside, e.g. the CLI.
func Validate(collection1, collection2, collection3 []int) error {
	for i := 1; i < len(collection1); i++ {
		if collection1[i] > collection1[i-1] {
			return fmt.Errorf("%w: collection_1 must be sorted from max to min, but %d comes after %d",
				ErrInvalidInput, collection1[i], collection1[i-1])
		}
	}
	ascending := []struct {
		name string
		c    []int
	}{{"collection_2", collection2}, {"collection_3", collection3}}
	for _, a := range ascending {
		name, c := a.name, a.c
		for i, v := range c {
			if v < 0 {
				return fmt.Errorf("%w: %s must start from min 0, but contains %d", ErrInvalidInput, name, v)
			}
			if i > 0 && v < c[i-1] {
				return fmt.Errorf("%w: %s must be sorted from min to max, but %d comes after %d",
					ErrInvalidInput, name, v, c[i-1])
			}
		}
	}
	return nil
}
