package ranger_test

import (
	"context"
	"testing"
	"time"

	"github.com/dcarbone/go-ranger"
)

type testItemT struct {
	Value   int
	Visited bool
}

func makeTestList(n int) []*testItemT {
	out := make([]*testItemT, n)
	for i := 0; i < n; i++ {
		out[i] = &testItemT{Value: i}
	}
	return out
}

func TestListRangerSync(t *testing.T) {
	t.Parallel()
	tl := makeTestList(50)

	lr := ranger.NewListRanger(tl...)

	lr.Range(func(i int, ts *testItemT) bool {
		ts.Visited = true
		return true
	})

	var visited int
	for _, itm := range tl {
		if itm.Visited {
			visited++
		}
	}
	if visited != len(tl) {
		t.Errorf("Expected to visit %d elements, only visited %d", len(tl), visited)
	}
}

func TestListRangerAsync(t *testing.T) {
	t.Parallel()

	t.Run("all", func(t *testing.T) {
		tl := makeTestList(50)

		lr := ranger.NewListRanger(tl...)

		<-lr.RangeAsync(context.Background(), func(_ context.Context, i int, ts *testItemT) {
			ts.Visited = true
		})

		var visited int
		for _, itm := range tl {
			if itm.Visited {
				visited++
			}
		}
		if visited != len(tl) {
			t.Errorf("Expected to visit %d elements, only visited %d", len(tl), visited)
		}
	})

	t.Run("chunked", func(t *testing.T) {
		tl := makeTestList(50000)

		lr := ranger.NewListRanger(tl...)

		<-lr.RangeAsyncChunked(context.Background(), 50, func(ctx context.Context, i int, ts *testItemT) {
			if ctx.Err() != nil {

			}
			ts.Visited = true
		})

		var visited int
		for _, itm := range tl {
			if itm.Visited {
				visited++
			}
		}
		if visited != len(tl) {
			t.Errorf("Expected to visit %d elements, only visited %d", len(tl), visited)
		}
	})

	t.Run("chunked-truncated", func(t *testing.T) {
		tl := makeTestList(50000)

		lr := ranger.NewListRanger(tl...)

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		<-lr.RangeAsyncChunked(ctx, 50, func(ctx context.Context, i int, ts *testItemT) {
			time.Sleep(50 * time.Millisecond)
			if ctx.Err() != nil {
				return
			}
			ts.Visited = true
		})

		var visited int
		for _, itm := range tl {
			if itm.Visited {
				visited++
			}
		}

		if len(tl) == visited {
			t.Error("Expected to NOT visit all elements")
		}
	})
}
