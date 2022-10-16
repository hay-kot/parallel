package parallel_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hay-kot/parallel"
)

var (
	ErrAbort = errors.New("abort")
)

// Test_Gather_Error ensures that the result of the errored task is not
// included in the gathered results.
func Test_Gather_Error(t *testing.T) {
	t.Parallel()

	producers := []parallel.ProducerFunc[string]{
		func(ctx context.Context) (string, error) {
			return "a", nil
		},
		func(ctx context.Context) (string, error) {
			return "invalid", ErrAbort
		},
		func(ctx context.Context) (string, error) {
			return "c", nil
		},
	}

	gathered, err := parallel.Gather(producers)

	if err != ErrAbort {
		t.Fatalf("expected error %v, got %v", ErrAbort, err)
	}

	for _, r := range gathered {
		if r == "invalid" {
			t.Fatalf("expected result to not contain %q", r)
		}
	}
}

// Test_Gather ensures that the results are gathered using the default
// options.
func Test_Gather(t *testing.T) {
	t.Parallel()

	results := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	producers := make([]parallel.ProducerFunc[string], 0, len(results))
	for _, r := range results {
		producers = append(producers, func(ctx context.Context) (string, error) {
			time.Sleep(100 * time.Millisecond)
			return r, nil
		})
	}

	gathered, err := parallel.Gather(producers)

	if err != nil {
		t.Fatal(err)
	}

	if len(gathered) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(gathered))
	}

	for i, r := range gathered {
		found := false
		for _, rr := range results {
			if r == rr {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("unexpected result %s at index %d", r, i)
		}
	}
}
