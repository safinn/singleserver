package singleserver

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// backupProgress draws an in-place progress bar for a backup. It is only active
// on an interactive terminal in text mode; for pipes, files, and --output json
// it is disabled and every method is a no-op, so machine output stays clean.
type backupProgress struct {
	w       io.Writer
	enabled bool

	mu         sync.Mutex
	label      string
	total      int64
	current    int64
	active     bool
	lastRender time.Time
}

// newBackupProgress enables the bar only when w is a text Output writing to an
// interactive terminal.
func newBackupProgress(w io.Writer) *backupProgress {
	if f, ok := progressTerminal(w); ok {
		return &backupProgress{w: f, enabled: true}
	}
	return &backupProgress{}
}

// progressTerminal returns the terminal file to draw on when w is a text Output
// connected to an interactive character device. NO_COLOR is intentionally not
// consulted: a plain bar is still useful without color.
func progressTerminal(w io.Writer) (*os.File, bool) {
	o, ok := w.(*Output)
	if !ok || o.json {
		return nil, false
	}
	f, ok := o.w.(*os.File)
	if !ok {
		return nil, false
	}
	if os.Getenv("TERM") == "dumb" {
		return nil, false
	}
	info, err := f.Stat()
	if err != nil || info.Mode()&os.ModeCharDevice == 0 {
		return nil, false
	}
	return f, true
}

// phase begins a new labeled bar with the given total in bytes.
func (p *backupProgress) phase(label string, total int64) {
	if !p.enabled {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.active {
		p.clearLocked()
	}
	p.label = label
	p.total = total
	p.current = 0
	p.active = true
	p.lastRender = time.Time{}
	p.renderLocked(true)
}

// add advances the current phase by n bytes.
func (p *backupProgress) add(n int64) {
	if !p.enabled || n == 0 {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.current += n
	p.renderLocked(false)
}

// set reports an absolute byte count for the current phase (used when progress
// is observed externally, e.g. by polling a file size).
func (p *backupProgress) set(n int64) {
	if !p.enabled {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.current = n
	p.renderLocked(false)
}

// finish erases the bar so the next output starts on a clean line.
func (p *backupProgress) finish() {
	if !p.enabled {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.active {
		p.clearLocked()
		p.active = false
	}
}

// renderLocked redraws the bar, throttled to ~10 fps unless forced or complete.
func (p *backupProgress) renderLocked(force bool) {
	now := time.Now()
	if !force && p.current < p.total && now.Sub(p.lastRender) < 100*time.Millisecond {
		return
	}
	p.lastRender = now
	fmt.Fprint(p.w, "\r"+p.barLocked()+"\x1b[K")
}

func (p *backupProgress) clearLocked() {
	fmt.Fprint(p.w, "\r\x1b[K")
}

func (p *backupProgress) barLocked() string {
	const width = 22
	ratio := 0.0
	if p.total > 0 {
		ratio = float64(p.current) / float64(p.total)
	}
	if ratio > 1 {
		ratio = 1
	}
	if ratio < 0 {
		ratio = 0
	}
	filled := int(ratio * width)
	bar := strings.Repeat("█", filled) + strings.Repeat(" ", width-filled)
	return fmt.Sprintf("%s  %s / %s  [%s] %3.0f%%",
		p.label, humanBytes(p.current), humanBytes(p.total), bar, ratio*100)
}

// progressReader counts bytes read and reports them to a backupProgress.
type progressReader struct {
	r io.Reader
	p *backupProgress
}

func (pr *progressReader) Read(b []byte) (int, error) {
	n, err := pr.r.Read(b)
	if n > 0 {
		pr.p.add(int64(n))
	}
	return n, err
}

func humanBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for x := n / unit; x >= unit; x /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGTPE"[exp])
}

// sqliteBackupTimeout scales the sqlite3 .backup timeout with database size so
// large databases do not trip a fixed deadline, with a generous floor.
func sqliteBackupTimeout(size int64) time.Duration {
	const floor = 60 * time.Second
	const bytesPerSec = 10 * 1024 * 1024 // conservative worst-case throughput
	scaled := time.Duration(size/bytesPerSec) * time.Second
	if scaled < floor {
		return floor
	}
	return scaled + floor
}
