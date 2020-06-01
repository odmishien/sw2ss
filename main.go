package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	start := time.Now()
	fmt.Printf("\x1b[36m%s\x1b[0m", "press Enter to stop your stopwatch!\n")
	bufio.NewScanner(os.Stdin).Scan()
	end := time.Now()
	diff := end.Sub(start)

	hours := int(diff.Hours()) % 24
	mins := int(diff.Minutes()) % 60
	secs := int(diff.Seconds()) % 60

	fmt.Printf("%d:%d:%d\n", hours, mins, secs)
}
