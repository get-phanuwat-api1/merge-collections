// Command merge runs merge.Merge from the command line.
//
//	go run ./cmd/merge -c1 9,7,5,3 -c2 0,2,4,8 -c3 1,6,10
//	[0 1 2 3 4 5 6 7 8 9 10]
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/get-phanuwat-api1/merge-collections"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("merge", flag.ContinueOnError)
	c1 := fs.String("c1", "", "collection_1, comma separated, sorted from max to min (e.g. 9,7,5)")
	c2 := fs.String("c2", "", "collection_2, comma separated, sorted from min (0) to max (e.g. 0,2,4)")
	c3 := fs.String("c3", "", "collection_3, comma separated, sorted from min (0) to max (e.g. 1,6,10)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	collection1, err := parseInts("c1", *c1)
	if err != nil {
		return err
	}
	collection2, err := parseInts("c2", *c2)
	if err != nil {
		return err
	}
	collection3, err := parseInts("c3", *c3)
	if err != nil {
		return err
	}
	if err := merge.Validate(collection1, collection2, collection3); err != nil {
		return err
	}

	_, err = fmt.Fprintln(out, merge.Merge(collection1, collection2, collection3))
	return err
}

func parseInts(flagName, s string) ([]int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return []int{}, nil
	}
	parts := strings.Split(s, ",")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		v, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil, fmt.Errorf("-%s: %q is not an integer", flagName, p)
		}
		out = append(out, v)
	}
	return out, nil
}
