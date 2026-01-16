package main

import (
	"fmt"
	"sort"
	"time"

	"kai/spinner"

	"github.com/fatih/color"
)

func main() {

	fmt.Println("Progress Bar Style Preview")
	fmt.Println("==========================")

	for name, style := range spinner.BarStyles {
		fmt.Printf("\nStyle: %s\n", name)

		bar := spinner.NewProgressBar(100).
			WithStyle(style).
			WithPrefix(fmt.Sprintf("%-15s ", name)).
			WithColor(color.New(color.FgCyan))

		bar.Start()

		for i := 0; i <= 100; i++ {
			bar.Update(i)
			time.Sleep(15 * time.Millisecond)
		}

		bar.Finish()
		time.Sleep(300 * time.Millisecond)
	}

	fmt.Println("\nAll bar styles previewed.")

	fmt.Println("Progress Spinner Style Preview")
	fmt.Println("==========================")

	ids := make([]int, 0, len(spinner.CharSets))
	for id := range spinner.CharSets {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	for _, id := range ids {
		fmt.Printf("Style ID: %d\n", id)

		s := spinner.NewSpinner().
			WithCharsetID(id).
			WithPrefix(fmt.Sprintf("[%02d] ", id)).
			WithColor(color.New(color.FgGreen)).
			WithSpeed(80 * time.Millisecond)

		s.Start()
		time.Sleep(800 * time.Millisecond)
		s.Stop()

		time.Sleep(120 * time.Millisecond)
	}

	fmt.Println("\n All spinner styles previewed.")

}
