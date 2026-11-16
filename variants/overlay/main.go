package main

import (
	"fmt"
	"os"

	"github.com/goome-ui/go-goome"
)

const appName = "My Goome App"

func main() {
	g, err := goome.Open(goome.Options{Title: appName})
	if err != nil {
		fmt.Fprintln(os.Stderr, "couldnt open an overlay:", err)
		os.Exit(1)
	}
	defer g.Close()
	g.Hotkey("alt+h", 1)
	ui, err := g.NewUI()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer ui.Free()
	win := goome.NewWindow(60, 60, 280, 0)
	clicks, hidden := 0, false
	for !win.Closed {
		for ev := range g.Events(16) {
			if ev.Type == goome.EventQuit {
				win.Closed = true
			}
			if ev.Type == goome.EventHotkey {
				hidden = !hidden
			}
			ui.Feed(ev)
		}
		if win.Closed {
			break
		}
		if ui.Begin(0) == nil {
			ui.End()
			continue
		}
		if !hidden {
			if ui.Window(appName, win, goome.WinMove|goome.WinClose) {
				ui.Label("im on top of everything")
				if ui.Button("Click me") {
					clicks++
				}
				ui.Label(fmt.Sprintf("clicked %d times", clicks))
				ui.Label("alt+h hides me")
			}
			ui.End()
		}
		ui.End()
	}
}
