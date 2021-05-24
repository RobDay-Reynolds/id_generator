package generator

import (
	"fmt"
	"testing"
	"time"
)

func TestGeneratorPerformance(t *testing.T) {
	generator := NewGenerator()

	startTime := time.Now()

	// Generate 100,000 IDs and see how long it takes
	for i := 0; i < 100000; i++ {
		generator.NextID()
	}

	endTime := time.Now()

	timeTaken := endTime.Sub(startTime)

	// If it takes more than a second then it wouldn't be performant enough
	if timeTaken.Seconds() > 1 {
		t.Errorf("Took too long to generate 100,000 IDs. Had to take less than one second, actually took %v", timeTaken)
	}

	fmt.Printf("Creating 100,000 IDs took %v\n", timeTaken)
}