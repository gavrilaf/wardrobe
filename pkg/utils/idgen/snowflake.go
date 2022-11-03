package idgen

import (
	"fmt"
	"sync"
	"time"
)

const (
	nodeIDBits   = 10
	sequenceBits = 12

	startEpoch = 1609459200000 // Custom Epoch 2021-01-01T00:00:00Z

	maxNodeID   = (1 << nodeIDBits) - 1
	maxSequence = (1 << sequenceBits) - 1
)

var SnowflakeTime = time.Now

// Snowflake unique IDs generator: https://en.wikipedia.org/wiki/Snowflake_ID
// The current implementation will generate unique IDs for 69 years since Midnight 1 Jan 2021
// Original implementation from Twitter (Scala): https://github.com/twitter-archive/snowflake/releases/tag/snowflake-2010
type Snowflake interface {
	NextID() (int64, error)
}

func NewSnowflake(nodeID int64) Snowflake {
	return &snowflake{nodeID: nodeID, lock: sync.Mutex{}}
}

//
type snowflake struct {
	nodeID   int64
	lastTm   int64
	sequence int64
	lock     sync.Mutex
}

func (s *snowflake) NextID() (int64, error) {
	tm := s.timestamp()

	if tm < s.lastTm {
		return 0, fmt.Errorf("clock moved backwards, refusing to generate id for %d milliseconds", s.lastTm-tm)
	}

	s.lock.Lock()

	if tm == s.lastTm {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			tm = s.waitNextMillis(tm)
		}
	} else {
		s.sequence = 0
	}

	s.lastTm = tm
	s.lock.Unlock()

	return tm<<(nodeIDBits+sequenceBits) | (s.nodeID << sequenceBits) | s.sequence, nil
}

// Block and wait till next millisecond
func (s *snowflake) waitNextMillis(tm int64) int64 {
	for tm == s.lastTm {
		tm = s.timestamp()
	}
	return tm
}

func (s *snowflake) timestamp() int64 {
	return SnowflakeTime().UnixMilli() - startEpoch
}
