package merge

import (
	"reflect"
	"testing"
)

func TestMergeSorted(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		left, right []int
		want        []int
	}{
		{name: "both empty", left: []int{}, right: []int{}, want: []int{}},
		{name: "left empty", left: nil, right: []int{1, 2}, want: []int{1, 2}},
		{name: "right empty", left: []int{1, 2}, right: nil, want: []int{1, 2}},
		{name: "interleaved", left: []int{1, 4, 7}, right: []int{2, 3, 8, 9}, want: []int{1, 2, 3, 4, 7, 8, 9}},
		{name: "left all smaller", left: []int{1, 2}, right: []int{5, 6}, want: []int{1, 2, 5, 6}},
		{name: "duplicates kept", left: []int{1, 3, 3}, right: []int{3, 4}, want: []int{1, 3, 3, 3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := mergeSorted(tt.left, tt.right); !reflect.DeepEqual(got, tt.want) {
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
