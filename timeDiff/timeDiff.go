package timeDiff

import (
	"fmt"
	"time"
)

// CalculateDiff calculates the elapsed time since a record was last checked
func CalculateDiff(timestamp time.Time) {
	timeDiff := time.Since(timestamp)
	fmt.Printf("Time since last check: %v\n", timeDiff)

	limit := time.Second * 43000

	if timeDiff > limit {
		fmt.Println("Time is greater than limit")
	} else {
		fmt.Println("Time OK")
	}
}
