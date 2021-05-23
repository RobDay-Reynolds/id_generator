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

	return generator
}

func (gen *Generator) NextID() uint64 {
	gen.mutex.Lock()
	defer gen.mutex.Unlock()

	unsignedTimestamp := uint64(timestamp())
	gen.sequence += 1

	return unsignedTimestamp<<(bitsForNodeID()+bitsForTimedRequests()) |
		gen.sequence<<bitsForNodeID() |
		gen.nodeID
}

// Placeholder to allow defining desired behavior and initial testing
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
