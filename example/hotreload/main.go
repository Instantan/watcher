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
