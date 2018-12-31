//+build js,wasm

package main

import (
	"time"

	// this helps with stutering, for example wjsu.Console.Log, 100% optional.
	. "github.com/OneOfOne/wjsu"
)

func main() {
	Console.Log("hello console", []int{1, -1}, map[string]float64{"a": 1, "z": -1})
	Document.QuerySelector("div#loader").SetInnerHTML("hello world")

	app := Document.QuerySelector("app")
	for t := range time.Tick(time.Second) {
		app.SetInnerHTML("time (UTC): " + t.UTC().String())
	}
}
