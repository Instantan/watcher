package main

import (
	"net/http"

	"github.com/Instantan/watcher"
)

func main() {
	watcher.HotReload()
	println("started.")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
