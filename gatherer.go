package parallel

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

type ProducerFunc[T any] func(ctx context.Context) (T, error)

// Gather[T any] uses the producer/consumer pattern to gather results from a list of producers.
// concurrently and return them as a slice of results.
//
// Supported options:
//   - WithMaxProcs
//   - WithContext
//
// Example:
//
//		results := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
//
//		producers := make([]parallel.ProducerFunc[string], 0, len(results))
//
//		for _, r := range results {
//			producers = append(producers, func(ctx context.Context) (string, error) {
//				time.Sleep(1 * time.Second)
//				return r, nil
//			})
//		}
//
//		gathered, err := parallel.Gather(producers)
//	    // Gathered order is _NOT_ guaranteed.
//	    // gathered == []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
func Gather[T any](tasks []ProducerFunc[T], opts ...optionsFunc) ([]T, error) {
	o := newOptions(opts...)

	arr := make([]T, 0, len(tasks))
	results := make(chan T)
	done := make(chan bool)

	var err error
	go func() {
		err = dispatch(o.ctx, tasks, done, o.maxProcs, results)
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range results {
			arr = append(arr, r)
		}
	}()

	<-done
	close(results)
	wg.Wait()

	return arr, err
}

func dispatch[T any](ctx context.Context, tasks []ProducerFunc[T], done chan bool, maxProcs int, result chan T) error {
	defer close(done)

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(maxProcs)

	for _, task := range tasks {
		task := task

		g.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				r, err := task(ctx)
				if err != nil {
					return err
				}
				result <- r
			}

			return nil
		})
	}

	return g.Wait()
}
