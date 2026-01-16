package main

import (
	"kai/bar"
	"kai/spinner"
	"os"
	"time"
)

func main() {
	s := spinner.New(os.Stdout, 100*time.Millisecond)
	s.Start()

	time.Sleep(2 * time.Second)
	s.Stop()

	bar := bar.New(100, 30, os.Stdout)
	for i := 0; i <= 100; i++ {
		bar.Render(i)
		time.Sleep(50 * time.Millisecond)
	}
	bar.Finish()
}
