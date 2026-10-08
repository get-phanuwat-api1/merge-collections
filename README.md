# merge-collections

A Go implementation of:

```
merge(int[] collection_1, int[] collection_2, int[] collection_3) : int[]
```

The function returns all elements sorted in **ascending** order, given that:

- `collection_1` is already sorted from **max to min**
- `collection_2` and `collection_3` are already sorted from **min (0) to max**

**No sort function is used.** The repository enforces this with `TestNoSortFunctionUsed`, which fails if any non-test file imports `sort` or `slices`.

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

Input that breaks these rules is rejected with a clear message:

```bash
go run ./cmd/merge -c1 1,2
# error: invalid input: collection_1 must be sorted from max to min, but 2 comes after 1
```

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
| `make fuzz` | fuzz test for 30s (`make fuzz FUZZTIME=2m`) |
| `make bench` | benchmark |
| `make run` | runs the example above |

---

## How it works

This is a **three-way merge** in a single pass:

```
collection_1 = [9, 7, 5, 3]   ← read from the END (smallest first):  3 → 5 → 7 → 9
collection_2 = [0, 2, 4, 8]   ← read from the START:                 0 → 2 → 4 → 8
collection_3 = [1, 6, 10]     ← read from the START:                 1 → 6 → 10

each step: compare the three current heads, append the smallest, advance that cursor
```

| | |
|---|---|
| Time | **O(n)**, where n is the total number of elements; each element is visited once |
| Extra space | O(n) for the result only (one allocation) |
| Inputs | never modified; the result is a new slice |
| Empty or nil inputs | return `[]`, never `nil` |
| Duplicates | kept |

`Merge` trusts its preconditions so that it stays a single O(n) pass. Call `Validate` first when the input comes from outside, as the CLI does.

## Tests

| Test | Covers |
|---|---|
| `TestMerge` | table-driven cases: interleaved values, empty or nil inputs, a single collection only, duplicates, all-equal values, non-overlapping ranges, different lengths, negatives in collection_1 |
| `TestMerge_DoesNotModifyInputs` / `ResultDoesNotAliasInputs` | inputs are untouched, and the result is a new slice |
| `TestMerge_LargeInput` | 300,000 elements |
| `TestValidate` | every precondition violation |
| `TestNoSortFunctionUsed` | the requirement that no sort function is used |
| `FuzzMerge` | random valid inputs: the result is ascending and contains exactly the input elements |
| `ExampleMerge` | runnable documentation example |
| `cmd/merge` `TestRun` | CLI parsing, output, and errors |

Coverage of the `merge` package is 100%.

## Project structure

```
merge-collections/
├── merge.go              # Merge + Validate
├── merge_test.go         # unit tests + benchmark + "no sort" guard
├── merge_fuzz_test.go    # fuzz test
├── example_test.go       # runnable example
├── cmd/merge/            # CLI
├── Makefile
└── go.mod                # no third-party dependencies
```
