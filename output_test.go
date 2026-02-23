package tinta

import (
	"testing"

	"github.com/varavelio/tinta/internal/assert"
)

func fakeEnv(m map[string]string) func(string) string {
	return func(key string) string { return m[key] }
}

// --- Color detection ---

func TestColorDetection(t *testing.T) {
	t.Run("NO_COLOR disables", func(t *testing.T) {
		assert.Equal(t, false, colorEnabled(fakeEnv(map[string]string{"NO_COLOR": "1"})))
	})

	t.Run("NO_COLORS disables", func(t *testing.T) {
		assert.Equal(t, false, colorEnabled(fakeEnv(map[string]string{"NO_COLORS": "1"})))
	})

	t.Run("DISABLE_COLORS disables", func(t *testing.T) {
		assert.Equal(t, false, colorEnabled(fakeEnv(map[string]string{"DISABLE_COLORS": "1"})))
	})

	t.Run("disable takes precedence over force", func(t *testing.T) {
		assert.Equal(t, false, colorEnabled(fakeEnv(map[string]string{
			"NO_COLOR":    "1",
			"FORCE_COLOR": "1",
		})))
	})

	t.Run("FORCE_COLOR enables", func(t *testing.T) {
		assert.Equal(t, true, colorEnabled(fakeEnv(map[string]string{"FORCE_COLOR": "1"})))
	})

	t.Run("CLICOLOR_FORCE enables", func(t *testing.T) {
		assert.Equal(t, true, colorEnabled(fakeEnv(map[string]string{"CLICOLOR_FORCE": "1"})))
	})

	t.Run("CLICOLOR=0 disables", func(t *testing.T) {
		assert.Equal(t, false, colorEnabled(fakeEnv(map[string]string{"CLICOLOR": "0"})))
	})

	t.Run("TERM=dumb disables", func(t *testing.T) {
		assert.Equal(t, false, colorEnabled(fakeEnv(map[string]string{"TERM": "dumb"})))
	})
}

func TestColorDetectionEdgeCases(t *testing.T) {
	t.Run("empty env falls through to terminal check", func(t *testing.T) {
		got := colorEnabled(fakeEnv(map[string]string{}))
		assert.Equal(t, false, got)
	})

	t.Run("TERM=dumb case insensitive", func(t *testing.T) {
		assert.Equal(t, false, colorEnabled(fakeEnv(map[string]string{"TERM": "DUMB"})))
		assert.Equal(t, false, colorEnabled(fakeEnv(map[string]string{"TERM": "Dumb"})))
	})

	t.Run("CLICOLOR=1 does not force enable", func(t *testing.T) {
		got := colorEnabled(fakeEnv(map[string]string{"CLICOLOR": "1"}))
		assert.Equal(t, false, got)
	})

	t.Run("FORCE_COLOR with any value enables", func(t *testing.T) {
		assert.Equal(t, true, colorEnabled(fakeEnv(map[string]string{"FORCE_COLOR": "0"})))
		assert.Equal(t, true, colorEnabled(fakeEnv(map[string]string{"FORCE_COLOR": "yes"})))
	})
}

// --- ANSI stripping ---

func TestStripANSI(t *testing.T) {
	t.Run("plain text unchanged", func(t *testing.T) {
		assert.Equal(t, "hello", stripANSI("hello"))
	})

	t.Run("strips single color", func(t *testing.T) {
		assert.Equal(t, "hello", stripANSI("\x1b[31mhello\x1b[0m"))
	})

	t.Run("strips multiple codes", func(t *testing.T) {
		assert.Equal(t, "hello", stripANSI("\x1b[37;44;1mhello\x1b[0m"))
	})

	t.Run("empty string", func(t *testing.T) {
		assert.Equal(t, "", stripANSI(""))
	})

	t.Run("only escape sequences", func(t *testing.T) {
		assert.Equal(t, "", stripANSI("\x1b[31m\x1b[0m"))
	})

	t.Run("mixed text and escapes", func(t *testing.T) {
		assert.Equal(t, "ab", stripANSI("\x1b[31ma\x1b[0m\x1b[32mb\x1b[0m"))
	})
}

func TestVisibleWidth(t *testing.T) {
	t.Run("plain text", func(t *testing.T) {
		assert.Equal(t, 5, visibleWidth("hello"))
	})

	t.Run("styled text", func(t *testing.T) {
		assert.Equal(t, 5, visibleWidth("\x1b[31mhello\x1b[0m"))
	})

	t.Run("empty string", func(t *testing.T) {
		assert.Equal(t, 0, visibleWidth(""))
	})

	t.Run("unicode characters", func(t *testing.T) {
		assert.Equal(t, 5, visibleWidth("hello"))
		assert.Equal(t, 4, visibleWidth("hola"))
	})
}
