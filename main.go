package main

import (
	"fmt"
	"os"

	"github.com/goome-ui/go-goome"
)

const appName = "My Goome App"

type state struct {
	clicks  int
	volume  float64
	overlay bool
	shown   bool
}

func openMode(overlay bool) (*goome.Goome, *goome.Window, error) {
	if !overlay {
		g, err := goome.OpenApp(appName, 640, 420)
		return g, nil, err
	}
	g, err := goome.Open(goome.Options{Title: appName})
	if err != nil {
		return nil, nil, err
	}
	x := 40.0
	if mons := g.Monitors(); len(mons) > 0 {
		x = float64(mons[0].Area.X + mons[0].Area.W - 340)
	}
	return g, goome.NewWindow(x, 60, 300, 0), nil
}

func contents(ui *goome.UI, s *state) {
	ui.Toggle("Overlay mode", &s.overlay)
	ui.Separator()
	text := uint32(0x00F2F4F8)
	if s.shown {
		text = 0xFFF2F4F8
	}
	ui.PushColor(goome.StyleText, ui.AnimColor("hello", text, 600, goome.EaseOut))
	s.shown = true
	ui.Heading("Hello from " + appName)
	ui.Pop(1)
	if ui.Button("Click me") {
		s.clicks++
	}
	ui.Label(fmt.Sprintf("clicked %d times", s.clicks))
	ui.SliderFmt("Volume", &s.volume, 0, 1, "%.2f")
}

func newUI(g *goome.Goome) *goome.UI {
	ui, err := g.NewUI()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return ui
}

func main() {
	s := &state{volume: 0.5}
	mode := false
	g, win, err := openMode(false)
	if err != nil {
		fmt.Fprintln(os.Stderr, "couldnt open a window:", err)
		os.Exit(1)
	}
	ui := newUI(g)
	for running := true; running; {
		for ev := range g.Events(16) {
			if ev.Type == goome.EventQuit || ev.Type == goome.EventClose {
				running = false
			}
			ui.Feed(ev)
		}
		if !running {
			break
		}
		if s.overlay != mode {
			ui.Free()
			g.Close()
			g, win, err = openMode(s.overlay)
			if err != nil {
				fmt.Fprintln(os.Stderr, "no overlay here, staying in window mode")
				s.overlay = false
				if g, win, err = openMode(false); err != nil {
					os.Exit(1)
				}
			}
			mode = s.overlay
			ui = newUI(g)
			if mode {
				fmt.Println("overlay mode")
			} else {
				fmt.Println("window mode")
			}
			continue
		}
		if ui.Begin(0) == nil {
			ui.End()
			continue
		}
		if mode {
			if ui.Window(appName, win, goome.WinMove|goome.WinCollapse) {
				contents(ui, s)
			}
			ui.End()
		} else {
			ui.App(appName, 0)
			contents(ui, s)
			ui.End()
		}
		ui.End()
	}
	ui.Free()
	g.Close()
}
