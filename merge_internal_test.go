package merge

import (
	"reflect"
	"testing"
)

func minToMax(x, y int) bool { return x < y }
func maxToMin(x, y int) bool { return x > y }

func TestMergeSorted(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		left, right []int
		less        func(x, y int) bool
		want        []int
	}{
		{name: "both empty", left: []int{}, right: []int{}, less: minToMax, want: []int{}},
		{name: "left empty", left: nil, right: []int{1, 2}, less: minToMax, want: []int{1, 2}},
		{name: "right empty", left: []int{1, 2}, right: nil, less: minToMax, want: []int{1, 2}},
		{name: "interleaved", left: []int{1, 4, 7}, right: []int{2, 3, 8, 9}, less: minToMax, want: []int{1, 2, 3, 4, 7, 8, 9}},
		{name: "left all smaller", left: []int{1, 2}, right: []int{5, 6}, less: minToMax, want: []int{1, 2, 5, 6}},
		{name: "duplicates kept", left: []int{1, 3, 3}, right: []int{3, 4}, less: minToMax, want: []int{1, 3, 3, 3, 4}},
		{name: "max to min order", left: []int{9, 5, 1}, right: []int{8, 2}, less: maxToMin, want: []int{9, 8, 5, 2, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := mergeSorted(tt.left, tt.right, tt.less); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeSorted(%v, %v) = %v, want %v", tt.left, tt.right, got, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in, want []int
	}{
		{in: []int{}, want: []int{}},
		{in: []int{1}, want: []int{1}},
		{in: []int{9, 7, 5, 3}, want: []int{3, 5, 7, 9}},
	}
	for _, tt := range tests {
		in := append([]int(nil), tt.in...)
		if got := reverse(tt.in); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("reverse(%v) = %v, want %v", tt.in, got, tt.want)
		}
		if !reflect.DeepEqual(tt.in, in) && len(in) > 0 {
			t.Errorf("reverse modified its input: %v", tt.in)
		}
	}
}
