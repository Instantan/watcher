# watcher
> File system watching and hotreloading

## Simple hot reloading for your go application
> main.go
```go
package main

import (
	"fmt"

	"github.com/Instantan/watcher"
)

func main() {
	watcher.HotReload()
	fmt.Printf("Started")
	<-make(chan struct{})
}
```