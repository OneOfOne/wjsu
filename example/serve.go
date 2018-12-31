//+build !js

package main

import (
	"net/http"
	"os/exec"
	"time"
)

func main() {
	time.AfterFunc(time.Second*1, func() {
		exec.Command("xdg-open", "http://localhost:8080").Start()
	})
	http.ListenAndServe(":8080", http.FileServer(http.Dir(".")))
}
