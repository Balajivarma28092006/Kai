package spinner

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/fatih/color"
)

type SpinnerStyle []string

type Spinner struct {
	frames     []string
	index      int
	prefix     string
	suffix     string
	color      *color.Color
	interval   time.Duration
	mu         sync.Mutex
	stop       chan struct{}
	running    bool
	HideCursor bool
}

func NewSpinner() *Spinner {
	return &Spinner{
		frames:     CharSets[43],
		interval:   100 * time.Millisecond,
		HideCursor: true,
		color:      color.New(color.FgCyan),
		stop:       make(chan struct{}),
	}
}

func (s *Spinner) WithCharsetID(id int) *Spinner {
	if style, ok := CharSets[id]; ok {
		s.mu.Lock()
		s.frames = style
		s.index = 0
		s.mu.Unlock()
	}
	return s
}

func (s *Spinner) WithFrames(frames []string) *Spinner {
	if len(frames) == 0 {
		return s
	}
	s.mu.Lock()
	s.frames = frames
	s.index = 0
	s.mu.Unlock()
	return s
}

func (s *Spinner) WithRandomStyle() *Spinner {
	keys := make([]int, 0, len(CharSets))
	for k := range CharSets {
		keys = append(keys, k)
	}

	randKey := keys[rand.Intn(len(keys))]
	return s.WithCharsetID(randKey)
}

var r *rand.Rand

func init() {
	// rand.Seed(time.Now().UnixNano())
	r = rand.New(rand.NewSource(time.Now().UnixNano()))

}

func ListStyles() {
	for id, frames := range CharSets {
		fmt.Printf("ID %-3d : %v\n", id, frames)
	}
}

func PreviewStyles() {
	for id, frames := range CharSets {
		fmt.Printf("ID %02d: ", id)
		for _, f := range frames[:min(6, len(frames))] {
			fmt.Print(f, " ")
		}
		fmt.Println()
	}
}

func (s *Spinner) WithPrefix(p string) *Spinner {
	s.prefix = p
	return s
}

func (s *Spinner) WithSuffix(suf string) *Spinner {
	s.suffix = suf
	return s
}

func (s *Spinner) WithColor(c *color.Color) *Spinner {
	s.color = c
	return s
}

func (s *Spinner) WithSpeed(d time.Duration) *Spinner {
	s.interval = d
	return s
}

func (s *Spinner) WithStyle(style SpinnerStyle) *Spinner {
	s.frames = style
	s.index = 0
	return s
}

func (s *Spinner) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}

	if s.HideCursor {
		// hides the cursor
		fmt.Print("\033[?25l")
	}
	s.running = true
	s.stop = make(chan struct{})
	s.mu.Unlock()

	ticker := time.NewTicker(s.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.render()
			case <-s.stop:
				ticker.Stop()
				fmt.Print("\r")
				return
			}
		}
	}()
}

func (s *Spinner) render() {
	s.mu.Lock()
	frame := s.frames[s.index]
	s.index = (s.index + 1) % len(s.frames)
	s.mu.Unlock()

	fmt.Printf("\r%s%s %s",
		s.prefix,
		s.color.Sprint(frame),
		s.suffix,
	)
}

func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}

	s.running = false
	close(s.stop)
	s.mu.Unlock()
	fmt.Println()
}
