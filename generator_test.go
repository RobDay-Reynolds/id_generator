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
		t.Errorf("Took too long to generate 100000 IDs. Had to take less than one second, actually took %v", timeTaken)
	}

	fmt.Printf("Creating 100000 IDs took %v\n", timeTaken)
}

func TestGeneratorInParallel(t *testing.T) {
	generator := NewGenerator()

	consumer := make(chan uint64)

	const idCount = 10000
	generateIDs := func() {
		for i:= 0; i < idCount; i++ {
			consumer <- generator.NextID()
		}
	}

	// Spin up five go funcs to create IDs in parallel
	const numThreads = 5
	for i := 0; i < numThreads; i++ {
		go generateIDs()
	}

	// Read off of the channel and check for duplicates
	idMap := map[uint64]bool{}
	for i := 0; i < idCount*numThreads; i++ {
		id := <-consumer
		if _, found := idMap[id]; found {
			t.Errorf("Duplicated ID")
		}
	}

	fmt.Printf("Created %v unique IDs\n", idCount*numThreads)
}