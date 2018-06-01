package main

import (
	"syscall/js"

	"lazyhackergo.com/browser"
)

func keepalive() {

	select {}
}

func main() {
	window := browser.Window()
}

// Slide movement
func HideHelpText(args []js.Value) {
	w := browser.Window()
	e := w.Document.GetElementById("help")
	s := browser.ElementStyle{
		property: "display",
		value:    "none",
	}
	e.SetStyle(s)
}

func GetSlideEl(args []js.Value) {

	w := browser.Window()
	slideEls := w.Document.QuerySelectorAll("section.slides > article")
	no := args[0]

}
