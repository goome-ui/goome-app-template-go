package main

import (
	"fmt"
	"os"
	"time"

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
	start := goome.NowMs()
	last := start
	fps, hidden := 0.0, false
	pace := goome.NewPacer(30)
	defer pace.Free()
	for running := true; running; {
		for ev := range g.Events(pace.Wait()) {
			if ev.Type == goome.EventQuit {
				running = false
			}
			if ev.Type == goome.EventHotkey {
				hidden = !hidden
			}
		}
		if !running {
			break
		}
		now := goome.NowMs()
		if now > last {
			fps = fps*0.9 + 100.0/float64(now-last)
		}
		last = now
		c := g.Begin(0)
		if c == nil {
			continue
		}
		if !hidden {
			c.RoundRect(20, 20, 160, 92, 10, 0xB0141820)
			c.Text(34, 30, time.Now().Format("15:04:05"), true, 0xFFFFFFFF)
			c.Text(34, 62, fmt.Sprintf("up %ds", (now-start)/1000), false, 0xFFC0C8D0)
			c.Text(34, 84, fmt.Sprintf("%.0f fps", fps), false, 0xFF6FDC5F)
		}
		g.End(c)
	}
}
