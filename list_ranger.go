package ranger

import (
	"context"
	"sync"
)

type ListRanger[T any] struct {
	list []T
}

func NewListRanger[T any](ts ...T) *ListRanger[T] {
	lr := ListRanger[T]{
		list: ts,
	}
	return &lr
}

// Len returns the number of configured
func (lr *ListRanger[T]) Len() int {
	return len(lr.list)
}

// Get will happily let you go beyond the length of the stored slice.  Don't do that.
func (lr *ListRanger[T]) Get(i int) T {
	return lr.list[i]
}

func (lr *ListRanger[T]) Range(cb func(int, T) bool) {
	for i, itm := range lr.list {
		if !cb(i, itm) {
			return
		}
	}
}

func (lr *ListRanger[T]) RangeAsync(ctx context.Context, cb func(context.Context, int, T)) <-chan struct{} {
	var (
		wg sync.WaitGroup

		done = make(chan struct{})
	)

	wg.Add(len(lr.list))

	for i := range lr.list {
		go func(i int) {
			defer wg.Done()
			cb(ctx, i, lr.list[i])
		}(i)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	return done
}

func (lr *ListRanger[T]) RangeAsyncChunked(ctx context.Context, n int, cb func(context.Context, int, T)) <-chan struct{} {
	var (
		wg sync.WaitGroup

		tickets = make(chan struct{}, n)
		done    = make(chan struct{})
	)

	for i := 0; i < n; i++ {
		tickets <- struct{}{}
	}

	wg.Add(len(lr.list))

	fn := func(i int) {
		defer wg.Done()
		cb(ctx, i, lr.list[i])
		tickets <- struct{}{}
	}

	go func() {
		for i := range lr.list {
			select {
			case <-tickets:
				go fn(i)
			case <-ctx.Done():
				wg.Done()
			}
		}
	}()

	go func() {
		wg.Wait()
		close(done)
		close(tickets)
		for range tickets {
		}
	}()

	return done
}

func RangeList[T any](cb func(int, T) bool, vs ...T) {
	(&ListRanger[T]{list: vs}).Range(cb)
}

func RangeListAsync[T any](ctx context.Context, cb func(context.Context, int, T), vs ...T) <-chan struct{} {
	return (&ListRanger[T]{list: vs}).RangeAsync(ctx, cb)
}

func RangeListAsyncChunked[T any](ctx context.Context, n int, cb func(context.Context, int, T), vs ...T) <-chan struct{} {
	return (&ListRanger[T]{list: vs}).RangeAsyncChunked(ctx, n, cb)
}
