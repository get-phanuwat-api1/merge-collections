// Package merge merges three pre-sorted integer collections into one ascending slice
// without using any sort function.
package merge

// Merge returns every element of the three collections in ascending order.
//
//	collection1: sorted from max to min   e.g. [9 7 5 3]
//	collection2: sorted from min to max   e.g. [0 2 4 8]
//	collection3: sorted from min to max   e.g. [1 6 10]
//	result:                               [0 1 2 3 4 5 6 7 8 9 10]
func Merge(collection1, collection2, collection3 []int) []int {
	// Step 1: collection1 is descending, so reverse it to make it ascending like the others.
	ascending1 := reverse(collection1)

	// Step 2: merge two sorted collections at a time.
	merged := mergeTwo(ascending1, collection2)
	return mergeTwo(merged, collection3)
}

// reverse returns a new slice with the elements of s in reverse order.
func reverse(s []int) []int {
	result := make([]int, 0, len(s))
	for i := len(s) - 1; i >= 0; i-- {
		result = append(result, s[i])
	}
	return result
}

// mergeTwo merges two ascending slices into one ascending slice.
//
// Look at the first remaining number of each slice, take the smaller one,
// and repeat. When one slice runs out, append the rest of the other.
func mergeTwo(a, b []int) []int {
	result := make([]int, 0, len(a)+len(b))
	i, j := 0, 0

	for i < len(a) && j < len(b) {
		if a[i] <= b[j] {
			result = append(result, a[i])
			i++
		} else {
			result = append(result, b[j])
			j++
		}
	}

	// One slice is finished; whatever is left in the other is already sorted.
	result = append(result, a[i:]...)
	result = append(result, b[j:]...)
	return result
}
