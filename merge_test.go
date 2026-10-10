package merge_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/get-phanuwat-api1/merge-collections"
)

func TestMerge(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                            string
		collection1, collection2, coll3 []int
		want                            []int
	}{
		{
			name:        "example: interleaved values",
			collection1: []int{9, 7, 5, 3},
			collection2: []int{0, 2, 4, 8},
			coll3:       []int{1, 6, 10},
			want:        []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:        "all empty returns empty",
			collection1: []int{}, collection2: []int{}, coll3: []int{},
			want: []int{},
		},
		{
			name: "all nil returns empty",
			want: []int{},
		},
		{
			name:        "only collection_1 (descending is reversed)",
			collection1: []int{5, 4, 3, 2, 1},
			want:        []int{1, 2, 3, 4, 5},
		},
		{
			name:        "only collection_2",
			collection2: []int{0, 1, 2},
			want:        []int{0, 1, 2},
		},
		{
			name:  "only collection_3",
			coll3: []int{0, 3, 9},
			want:  []int{0, 3, 9},
		},
		{
			name:        "duplicates inside and across collections are kept",
			collection1: []int{5, 5, 2, 0},
			collection2: []int{0, 2, 2, 5},
			coll3:       []int{0, 5},
			want:        []int{0, 0, 0, 2, 2, 2, 5, 5, 5, 5},
		},
		{
			name:        "all values equal",
			collection1: []int{3, 3},
			collection2: []int{3},
			coll3:       []int{3, 3, 3},
			want:        []int{3, 3, 3, 3, 3, 3},
		},
		{
			name:        "collections do not overlap",
			collection1: []int{30, 20, 10},
			collection2: []int{0, 1, 2},
			coll3:       []int{100, 200},
			want:        []int{0, 1, 2, 10, 20, 30, 100, 200},
		},
		{
			name:        "different lengths",
			collection1: []int{1},
			collection2: []int{0, 2, 3, 4, 5, 6, 7},
			coll3:       []int{},
			want:        []int{0, 1, 2, 3, 4, 5, 6, 7},
		},
		{
			name:        "collection_1 may contain negatives",
			collection1: []int{4, -1, -7},
			collection2: []int{0, 4},
			coll3:       []int{1},
			want:        []int{-7, -1, 0, 1, 4, 4},
		},
		{
			name:        "single elements",
			collection1: []int{2},
			collection2: []int{1},
			coll3:       []int{0},
			want:        []int{0, 1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := merge.Merge(tt.collection1, tt.collection2, tt.coll3)
			if len(got) == 0 && len(tt.want) == 0 {
				return // nil and [] are both an empty result
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge(%v, %v, %v)\n got  %v\n want %v", tt.collection1, tt.collection2, tt.coll3, got, tt.want)
			}
		})
	}
}

func TestMerge_DoesNotModifyInputs(t *testing.T) {
	t.Parallel()
	c1, c2, c3 := []int{9, 5, 1}, []int{0, 4}, []int{2, 8}
	merge.Merge(c1, c2, c3)

	if !reflect.DeepEqual(c1, []int{9, 5, 1}) || !reflect.DeepEqual(c2, []int{0, 4}) || !reflect.DeepEqual(c3, []int{2, 8}) {
		t.Errorf("inputs were modified: %v %v %v", c1, c2, c3)
	}
}

func TestMerge_ResultDoesNotAliasInputs(t *testing.T) {
	t.Parallel()
	c2 := []int{0, 1, 2}
	got := merge.Merge(nil, c2, nil)
	got[0] = 99
	if c2[0] != 0 {
		t.Error("changing the result changed collection_2: result must be a new slice")
	}
}

func TestMerge_LargeInput(t *testing.T) {
	t.Parallel()
	const n = 100_000 // 300,000 elements in total
	c1, c2, c3 := make([]int, n), make([]int, n), make([]int, n)
	for i := range n {
		c1[i] = 3*(n-1-i) + 0 // 3n-3 ... 3, 0   (descending)
		c2[i] = 3*i + 1       // 1, 4, 7 ...      (ascending)
		c3[i] = 3*i + 2       // 2, 5, 8 ...      (ascending)
	}

	got := merge.Merge(c1, c2, c3)

	if len(got) != 3*n {
		t.Fatalf("len = %d, want %d", len(got), 3*n)
	}
	for i, v := range got {
		if v != i {
			t.Fatalf("got[%d] = %d, want %d", i, v, i)
		}
	}
}

// TestNoSortFunctionUsed enforces the requirement "Don't allow to use any sort function":
// non-test source files must not use the sort package or slices.Sort* functions.
func TestNoSortFunctionUsed(t *testing.T) {
	t.Parallel()
	bannedSlicesFuncs := map[string]bool{"Sort": true, "SortFunc": true, "SortStableFunc": true}

	err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			return err
		}
		for _, imp := range file.Imports {
			if pkg, _ := strconv.Unquote(imp.Path.Value); pkg == "sort" {
				t.Errorf("%s imports %q: sort functions are not allowed", path, pkg)
			}
		}
		ast.Inspect(file, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "slices" && bannedSlicesFuncs[sel.Sel.Name] {
				t.Errorf("%s uses slices.%s: sort functions are not allowed", path, sel.Sel.Name)
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkMerge(b *testing.B) {
	const n = 10_000
	c1, c2, c3 := make([]int, n), make([]int, n), make([]int, n)
	for i := range n {
		c1[i], c2[i], c3[i] = 3*(n-i), 3*i+1, 3*i+2
	}
	b.ResetTimer()
	for range b.N {
		merge.Merge(c1, c2, c3)
	}
}
