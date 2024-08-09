# go-ranger
Collection of little utilities to assist with ranging over slices in go

<!-- TOC -->
* [go-ranger](#go-ranger)
  * [Examples](#examples)
  * [ListRanger](#listranger)
<!-- TOC -->

## Examples

Any / all examples are located in the [examples](./examples) directory.

## ListRanger

To range over a slice of items, you can use the [ListRanger](./list_ranger.go) type:

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dcarbone/go-ranger"
)

func main() {
	type item struct {
		Value int
		Visited bool
    }

	items := make([]*item, 50)
	for i := 0; i < 50; i++ {
		items[i] = &item{Value: i}
    }

	// range synchronously
	ranger.RangeList(func(i int, itm *item) bool {
		itm.Visited = true
		fmt.Printf("item %d visited synchronously\n", itm.Value)
        return true
	}, items...)
	
	// range asynchronously
	<-ranger.RangeListAsync(context.Background(), func(ctx context.Context, i int, itm *item) {
		itm.Visited = true
		fmt.Printf("item %d visited asynchronously\n", itm.Value)
    }, items...)
	
	// range asynchronously in chunks
	
	<-ranger.RangeListAsyncChunked(context.Background(), 10, func(ctx context.Context, i int, itm *item) {
		itm.Visited = true
		fmt.Printf("item %d visited asynchronously\n", itm.Value)
		time.Sleep(500 * time.Second)
    }, items...)
}
```