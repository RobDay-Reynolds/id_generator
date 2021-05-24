package generator

import (
	"math"
	"sync"
	"time"
)

const (
	MaxNumNodes = 1024
	MaxRequestsPerSecond = 100000
)

// Generator is a distributed ID generator.
type Generator struct {
	mutex    *sync.Mutex
	sequence uint64
	nodeID   uint64
}

func NewGenerator() *Generator {
	generator := new(Generator)
	generator.mutex = new(sync.Mutex)
	generator.nodeID = uint64(nodeID())

	// Theoretically, NextID could fail to return a unique ID if the node crashed
	// and restarted in the same millisecond. Similarly, if the entire system crashed
	// and restarted in the same millisecond each node could create duplicate IDs.
	// By sleeping for a single millisecond at startup we remove that as a possibility.
	time.Sleep(1 * time.Millisecond)

	return generator
}

func (gen *Generator) NextID() uint64 {
	// By including a mutex lock/unlock, we ensure that no two local calls
	// can process simultaneously
	gen.mutex.Lock()
	defer gen.mutex.Unlock()

	unsignedTimestamp := uint64(timestamp())

	// Increment our sequence to ensure we get a different ID than the last
	// local call, even if millisecond timestamp is the same
	gen.sequence += 1

	// If our sequence has reached our max throughput per second, reset it
	// back to 0. This is actually ~1000 times less frequently than we could
	// probably cycle this, but there's no way of knowing if the 100,000
	// requests per second all happen in the same millisecond
	if gen.sequence > MaxRequestsPerSecond {
		gen.sequence = 0
	}

	// By building an unsigned int from a millisecond timestamp, the node id,
	// and a monotonically increasing sequence, we ensure that no two calls will
	// return the same value
	return unsignedTimestamp<<(bitsForNodeID()+bitsForTimedRequests()) |
		gen.sequence<<bitsForNodeID() |
		gen.nodeID
}

// Returns millisecond accurate epoch time
func timestamp() int64 {
	return time.Now().Unix()
}

// Placeholder to allow defining desired behavior and initial testing
// challenge describes node_id as defined elsewhere
func nodeID() int {
	return 1
}

func bitsForNodeID() int {
	return int(math.Ceil(math.Log2(MaxNumNodes)))
}

func bitsForTimedRequests() int {
	return int(math.Ceil(math.Log2(MaxRequestsPerSecond)))
}
