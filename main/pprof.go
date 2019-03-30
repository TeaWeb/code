package main

import (
	"github.com/TeaWeb/code/teaweb"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		http.ListenAndServe("127.0.0.1:6060", nil)
	}()
	teaweb.Start()
}
