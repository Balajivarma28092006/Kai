package main

import (
	"time"

	"kai/spinner"

	"github.com/fatih/color"
)

func main() {
	bar := spinner.NewProgressBar(100).
		WithStyle(spinner.BarArrowAscii).
		WithPrefix("Doing Something").
		WithColor(color.New(color.FgCyan))

	bar.Start()

	for i := 0; i <= 100; i++ {
		bar.Update(i)
		time.Sleep(15 * time.Millisecond)
	}

	bar.Finish()
	time.Sleep(300 * time.Millisecond)

	s := spinner.NewSpinner().
		WithCharsetID(52).
		WithColor(color.New(color.FgGreen)).
		WithSpeed(100 * time.Millisecond)

	s.Start()
	for i := 0; i < 10; i++ {
		time.Sleep(5000 * time.Millisecond)
	}
	s.Stop()

	time.Sleep(120 * time.Millisecond)
}
