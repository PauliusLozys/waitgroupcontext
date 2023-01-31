# WaitGroupContext package
This package adds `WaitGroupContext` similar to `sync.WaitGroup` but allows to pass `context.Context` to cancel from `Wait()` method early. 
```sh
go get github.com/PauliusLozys/waitcroupcontext
```

Usage
```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/PauliusLozys/waitgroupcontext"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wgc := waitgroupcontext.New(ctx)

	wgc.Add(1)
	go func() {
		defer wgc.Sub()
		// Work...
	}()

	go func() {
		select {
		case <-time.After(time.Second):
			cancel() // cancel() WaitGroupContext's ctx which will cause wgc.Wait() to finish
		case <-wgc.Done():
			fmt.Println("This will be printed when wgc.Wait() function finishes")
		}
	}()

	wgc.Wait()
}
```