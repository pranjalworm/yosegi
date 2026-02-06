package utils

import (
	"fmt"
	"time"
)

func TrackTime(name string) func() {

	start := time.Now()

	return func() {
		fmt.Printf("%s took %s\n", name, time.Since(start))
	}
}
