package merge_test

import (
	"fmt"

	"github.com/get-phanuwat-api1/merge-collections"
)

func ExampleMerge() {
	collection1 := []int{9, 7, 5, 3} // sorted from max to min
	collection2 := []int{0, 2, 4, 8} // sorted from min (0) to max
	collection3 := []int{1, 6, 10}   // sorted from min (0) to max

	fmt.Println(merge.Merge(collection1, collection2, collection3))
	// Output: [0 1 2 3 4 5 6 7 8 9 10]
}
