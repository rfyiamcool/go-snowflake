package snowflake

import (
	"errors"
	"hash/crc32"
	"math/rand"
	"net"
	"sync"
	"time"
)

const (
	nano               = 1000 * 1000
	WorkerIDBits       = 10              // worker id
	MaxWorkerID  int64 = -1 ^ (-1 << 10) // worker id mask
	MaxSequence        = -1 ^ (-1 << 12) // sequence mask
	SequenceBits       = 12              // sequence
)

var (
	sf = NewSnowFlake(GetDefaultWorkID())

	Since = time.Date(2012, 1, 0, 0, 0, 0, 0, time.UTC).UnixNano() / nano
)

type SnowFlake struct {
	lastTimestamp int64
	workerID      uint32
	sequence      uint32
	lock          sync.Mutex
}

func Init(id int64) {
	if id == 0 {
		id = GetDefaultWorkID()
	}
	sf = NewSnowFlake(id)
}

func NewSnowFlake(id int64) *SnowFlake {
	workerID := id
	if id < 0 || id > MaxWorkerID {
		workerID = (workerID%MaxWorkerID + MaxWorkerID) % MaxWorkerID
	}

	return &SnowFlake{workerID: uint32(workerID)}
}

func (sf *SnowFlake) int64() int64 {
	return (sf.lastTimestamp << (WorkerIDBits + SequenceBits)) |
		(int64(sf.workerID) << SequenceBits) |
		int64(sf.sequence)
}

// Next get id
func (sf *SnowFlake) Next() (int64, error) {
	sf.lock.Lock()
	defer sf.lock.Unlock()

	ts := timestamp()
	if ts < sf.lastTimestamp {
		// avoid cause timestamp rollback, try to get timestamp greater than lastTimestamp
		ts = waitUntilMillis(sf.lastTimestamp)
		if ts < sf.lastTimestamp {
			return 0, errors.New("timestamp cause rollback")
		}
	}

	if ts == sf.lastTimestamp {
		sf.sequence = (sf.sequence + 1) & MaxSequence
	} else {
		sf.sequence = 0
	}

	sf.lastTimestamp = ts
	return sf.int64(), nil
}

// GetWorkerID get current worker ID
func (sf *SnowFlake) GetWorkerID() uint32 {
	return sf.workerID
}

// GetTimeFromID extracts timestamp from an existing ID, unit Millisecond
func (sf *SnowFlake) GetTimeFromID(id int64) int64 {
	return id>>(WorkerIDBits+SequenceBits) + Since
}

// waitUntilMillis
func waitUntilMillis(ts int64) int64 {
	cur := timestamp()
	for i := 0; i < 100000; i++ { // 10s
		if cur > ts {
			return cur
		}

		cur = timestamp()
		time.Sleep(time.Duration(50 * time.Microsecond)) // delay 50 us
	}

	return cur
}

// timestamp relative timestamp
func timestamp() int64 {
	return int64(time.Now().UnixNano()/nano - Since)
}

// GetDefaultWorkID
func GetDefaultWorkID() int64 {
	var id int64

	ift, err := net.Interfaces()
	if err != nil {
		rand.Seed(time.Now().UnixNano())
		id = int64(rand.Uint32()) % MaxWorkerID
	} else {
		h := crc32.NewIEEE()
		for _, value := range ift {
			h.Write(value.HardwareAddr)
		}
		id = int64(h.Sum32()) % MaxWorkerID
	}

	return id & MaxWorkerID
}

func Next() (int64, error) {
	return sf.Next()
}

func GetTimeFromID(id int64) int64 {
	return sf.GetTimeFromID(id)
}

func GetWorkerID() int64 {
	return int64(sf.GetWorkerID())
}
