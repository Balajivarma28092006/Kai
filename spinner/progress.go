package spinner

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type BarStyle struct {
	Fill  string
	Empty string
}

type ProgressBar struct {
	total       int
	current     int
	width       int
	fillChar    string
	emptyChar   string
	prefix      string
	suffix      string
	showPercent bool
	color       *color.Color
	mu          sync.Mutex
	done        chan struct{}
	running     bool
}

func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{
		total:       total,
		current:     0,
		width:       40,
		fillChar:    BarLineMinimal.Fill,
		emptyChar:   BarLineMinimal.Empty,
		showPercent: true,
		color:       color.New(color.FgCyan),
		done:        make(chan struct{}),
	}
}

func (pb *ProgressBar) WithWidth(width int) *ProgressBar {
	pb.width = width
	return pb
}

func (pb *ProgressBar) WithPrefix(prefix string) *ProgressBar {
	pb.prefix = prefix
	return pb
}

func (pb *ProgressBar) WithSuffix(suffix string) *ProgressBar {
	pb.suffix = suffix
	return pb
}

func (pb *ProgressBar) WithColor(c *color.Color) *ProgressBar {
	pb.color = c
	return pb
}

func (pb *ProgressBar) WithStyle(style BarStyle) *ProgressBar {
	pb.fillChar = style.Fill
	pb.emptyChar = style.Empty
	return pb
}

func (pb *ProgressBar) WithChars(fill, empty string) *ProgressBar {
	pb.fillChar = fill
	pb.emptyChar = empty
	return pb
}

func (pb *ProgressBar) Start() {
	pb.mu.Lock()
	pb.running = true
	pb.mu.Unlock()

	go pb.renderer()
}

func (pb *ProgressBar) renderer() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-pb.done:
			pb.render()
			fmt.Println()
			return
		case <-ticker.C:
			pb.render()
		}
	}
}

func (pb *ProgressBar) Update(current int) {
	pb.mu.Lock()
	pb.current = current
	pb.mu.Unlock()
}

func (pb *ProgressBar) Increment() {
	pb.mu.Lock()
	pb.current++
	if pb.current > pb.total {
		pb.current = pb.total
	}
	pb.mu.Unlock()
}

func (pb *ProgressBar) render() {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	percent := float64(pb.current) / float64(pb.total) * 100
	filled := int(float64(pb.width) * float64(pb.current) / float64(pb.total))

	bar := strings.Repeat(pb.fillChar, filled) + strings.Repeat(pb.emptyChar, pb.width-filled)
	coloredBar := pb.color.Sprint(bar)

	output := fmt.Sprintf("\r%s[%s] %.0f%% (%d/%d) %s",
		pb.prefix,
		coloredBar,
		percent,
		pb.current,
		pb.total,
		pb.suffix,
	)

	fmt.Print(output)
}

func (pb *ProgressBar) Finish() {
	pb.mu.Lock()
	pb.current = pb.total
	pb.running = false
	pb.mu.Unlock()
	close(pb.done)
}
