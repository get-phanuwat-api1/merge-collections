package merge

import (
	"errors"
	"fmt"
)

// ErrInvalidInput is wrapped by Validate when a collection breaks the preconditions.
var ErrInvalidInput = errors.New("invalid input")

// Validate checks the preconditions of Merge. Merge trusts its input;
// call Validate first when the input comes from outside (e.g. the CLI).
func Validate(collection1, collection2, collection3 []int) error {
	for i := 1; i < len(collection1); i++ {
		if collection1[i] > collection1[i-1] {
			return fmt.Errorf("%w: collection_1 must be sorted from max to min, but %d comes after %d",
				ErrInvalidInput, collection1[i], collection1[i-1])
		}
	}
	if err := validateAscendingFromZero("collection_2", collection2); err != nil {
		return err
	}
	return validateAscendingFromZero("collection_3", collection3)
}

func validateAscendingFromZero(name string, c []int) error {
	for i, v := range c {
		if v < 0 {
			return fmt.Errorf("%w: %s must start from min 0, but contains %d", ErrInvalidInput, name, v)
		}
		if i > 0 && v < c[i-1] {
			return fmt.Errorf("%w: %s must be sorted from min to max, but %d comes after %d",
				ErrInvalidInput, name, v, c[i-1])
		}
	}
	return nil
}
