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
	g.Hotkey("alt+m", 1)
	ui, err := g.NewUI()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer ui.Free()
	menu := &goome.Window{X: 80, Y: 80, W: 300, Closed: true}
	fly, god, speed := false, false, 1.0
	for running := true; running; {
		for ev := range g.Events(16) {
			if ev.Type == goome.EventQuit {
				running = false
			}
			if ev.Type == goome.EventHotkey {
				menu.Closed = !menu.Closed
			}
			ui.Feed(ev)
		}
		if !running {
			break
		}
		c := ui.Begin(0)
		if c == nil {
			ui.End()
			continue
		}
		if !menu.Closed {
			if ui.Window(appName, menu, goome.WinMove|goome.WinClose) {
				ui.Toggle("Fly", &fly)
				ui.Checkbox("God mode", &god)
				ui.SliderFmt("Speed", &speed, 1, 10, "%.1f")
				if ui.Button("Respawn") {
					ui.Toast("respawned", 2000)
				}
			}
			ui.End()
		} else {
			c.RoundRect(20, 20, 170, 30, 8, 0xC0202838)
			c.Text(32, 26, "alt+m opens the menu", false, 0xFFFFFFFF)
		}
		ui.End()
	}
}
