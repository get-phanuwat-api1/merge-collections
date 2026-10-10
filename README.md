# merge-collections

A Go implementation of:

```
merge(int[] collection_1, int[] collection_2, int[] collection_3) : int[]
```

The function returns all elements sorted in **ascending** order, given that:

- `collection_1` is already sorted from **max to min**
- `collection_2` and `collection_3` are already sorted from **min (0) to max**

**No sort function is used.** The repository enforces this with `TestNoSortFunctionUsed`, which fails if any non-test file imports the `sort` package or calls `slices.Sort`, `slices.SortFunc`, or `slices.SortStableFunc`. The input is already sorted, so the implementation only merges it.

---

## Setup

The only dependency is **Go 1.22 or newer**. The project uses the Go standard library only, so there is nothing else to install.

```bash
go version   # check Go is installed: https://go.dev/dl/
```

```bash
git clone https://github.com/get-phanuwat-api1/merge-collections.git
cd merge-collections
```

## Run the code

```bash
go run ./cmd/merge -c1 9,7,5,3 -c2 0,2,4,8 -c3 1,6,10
# [0 1 2 3 4 5 6 7 8 9 10]
```

| Flag | Meaning |
|---|---|
| `-c1` | collection_1, comma separated, sorted from max to min |
| `-c2` | collection_2, comma separated, sorted from min (0) to max |
| `-c3` | collection_3, comma separated, sorted from min (0) to max |

The input must follow the rules in the table above, which are the "given" of the specification. The function relies on them instead of sorting.

To use it as a library:

```go
import "github.com/get-phanuwat-api1/merge-collections"

result := merge.Merge([]int{9, 7, 5, 3}, []int{0, 2, 4, 8}, []int{1, 6, 10})
// [0 1 2 3 4 5 6 7 8 9 10]
```

## Run the unit tests

```bash
go test ./...            # all unit tests
```

```bash
go test -v -cover ./...  # verbose, with coverage
```

Or with `make`:

| Command | What it does |
|---|---|
| `make test` | unit tests with the race detector |
| `make cover` | tests with coverage |
| `make bench` | benchmark |
| `make run` | runs the example above |

---

## How it works

Every collection is **already sorted**, so nothing needs to be sorted again. The function only merges:

```go
func Merge(collection1, collection2, collection3 []int) []int {
	ascending1 := reverse(collection1)             // 1. max→min becomes min→max
	merged := mergeSorted(ascending1, collection2) // 2. merge two sorted collections
	return mergeSorted(merged, collection3)        // 3. merge in the third one
}
```

```
reverse([9, 7, 5, 3])                               → [3, 5, 7, 9]
mergeSorted([3, 5, 7, 9], [0, 2, 4, 8])             → [0, 2, 3, 4, 5, 7, 8, 9]
mergeSorted([0, 2, 3, 4, 5, 7, 8, 9], [1, 6, 10])   → [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
```

`mergeSorted(left, right)` works like merging two sorted piles of cards. It looks at the top card of each pile, takes the smaller one, and repeats. When one pile runs out, the rest of the other pile is already in order.

| | |
|---|---|
| Time | **O(n)**, where n is the total number of elements; each element is visited a constant number of times, with no sorting |
| Inputs | never modified; the result is a new slice |
| Empty input | returns an empty slice |
| Duplicates | kept |

## Tests

| Test | Covers |
|---|---|
| `TestMerge` | table-driven cases: interleaved values, empty or nil inputs, a single collection only, duplicates, all-equal values, non-overlapping ranges, different lengths, negatives in collection_1 |
| `TestMerge_DoesNotModifyInputs` / `ResultDoesNotAliasInputs` | inputs are untouched, and the result is a new slice |
| `TestMerge_LargeInput` | 300,000 elements |
| `TestMergeSorted` / `TestReverse` | the two helpers on their own |
| `TestNoSortFunctionUsed` | the requirement that no sort function is used |
| `ExampleMerge` | runnable documentation example |
| `cmd/merge` `TestRun` | CLI parsing, output, and errors |

Coverage of the `merge` package is 100%.

## Project structure

```
merge-collections/
├── merge.go              # Merge, mergeSorted, reverse
├── merge_internal_test.go # tests for mergeSorted and reverse
├── merge_test.go         # unit tests + benchmark + "no sort" guard
├── example_test.go       # runnable example
├── cmd/merge/            # CLI
├── Makefile
└── go.mod                # no third-party dependencies
```
