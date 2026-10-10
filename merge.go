// Package merge merges three pre-sorted integer collections into one ascending slice
// without using any sort function.
package merge

// Merge returns every element of the three collections in ascending order.
//
// Given (from the specification):
//
//	collection1: already sorted from max to min   e.g. [9 7 5 3]
//	collection2: already sorted from min to max   e.g. [0 2 4 8]
//	collection3: already sorted from min to max   e.g. [1 6 10]
//
// Because every collection is already sorted, nothing needs to be sorted again:
// flip collection1 to ascending, then merge the sorted collections two at a time.
func Merge(collection1, collection2, collection3 []int) []int {
	minToMax := func(x, y int) bool { return x < y }

	ascending1 := reverse(collection1)
	merged := mergeSorted(ascending1, collection2, minToMax)
	return mergeSorted(merged, collection3, minToMax)
}

// mergeSorted merges two slices that are already sorted by less into one sorted slice.
// less decides the order: x < y for min to max, x > y for max to min.
//
// Compare the first remaining number of each slice, take the one that comes first,
// and repeat. When one slice runs out, the rest of the other is already in order.
func mergeSorted(left, right []int, less func(x, y int) bool) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if less(right[j], left[i]) {
			result = append(result, right[j])
			j++
		} else {
			result = append(result, left[i])
			i++
		}
	}
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

// reverse returns a new slice with the elements of s in reverse order.
func reverse(s []int) []int {
	result := make([]int, 0, len(s))
	for i := len(s) - 1; i >= 0; i-- {
		result = append(result, s[i])
	}
	return result
}
