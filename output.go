package tinta

import (
	"io"
	"os"
	"strings"
	"sync"
)

var (
	mu      sync.RWMutex
	output  io.Writer = os.Stdout
	enabled           = detectColor()
)

// SetOutput changes the default writer used by Print, Println and Printf
// on both [TextStyle] and [BoxStyle]. It is safe for concurrent use.
func SetOutput(w io.Writer) {
	mu.Lock()
	output = w
	mu.Unlock()
}

// ForceColors overrides automatic color detection. It is safe for concurrent
// use.
func ForceColors(on bool) {
	mu.Lock()
	enabled = on
	mu.Unlock()
}

func getOutput() io.Writer {
	mu.RLock()
	w := output
	mu.RUnlock()
	return w
}

func isEnabled() bool {
	mu.RLock()
	v := enabled
	mu.RUnlock()
	return v
}

func detectColor() bool {
	return colorEnabled(os.Getenv)
}

func colorEnabled(getenv func(string) string) bool {
	if getenv("NO_COLOR") != "" || getenv("NO_COLORS") != "" || getenv("DISABLE_COLORS") != "" {
		return false
	}
	if getenv("FORCE_COLOR") != "" || getenv("CLICOLOR_FORCE") != "" {
		return true
	}
	if getenv("CLICOLOR") == "0" {
		return false
	}
	if strings.EqualFold(getenv("TERM"), "dumb") {
		return false
	}
	return isTerminal(os.Stdout)
}

func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok || f == nil {
		return false
	}
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return stat.Mode()&os.ModeCharDevice != 0
}

func stripANSI(s string) string {
	if strings.IndexByte(s, '\x1b') < 0 {
		return s
	}

	var b strings.Builder
	b.Grow(len(s))

	i := 0
	for i < len(s) {
		if s[i] != '\x1b' {
			b.WriteByte(s[i])
			i++
			continue
		}
		if i+1 >= len(s) {
			break
		}
		switch s[i+1] {
		case '[':
			j := i + 2
			for j < len(s) && s[j] < 0x40 || s[j] > 0x7E {
				if s[j] >= 0x40 && s[j] <= 0x7E {
					break
				}
				j++
			}
			if j < len(s) {
				j++
			}
			i = j
		case ']':
			j := i + 2
			for j < len(s) {
				if s[j] == '\x07' {
					j++
					break
				}
				if s[j] == '\x1b' && j+1 < len(s) && s[j+1] == '\\' {
					j += 2
					break
				}
				j++
			}
			i = j
		default:
			i += 2
		}
	}
	return b.String()
}

func visibleWidth(s string) int {
	stripped := stripANSI(s)
	n := 0
	for range stripped {
		n++
	}
	return n
}
