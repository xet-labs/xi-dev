package service

import (
	"runtime"
	"time"

	"github.com/rs/zerolog/log"
)

type StatsService struct{}

var Stats = &StatsService{}

// MemD runs memory stats logging as a daemon at an interval (in seconds).
func (s *StatsService) MemD(intervalSec int) {
	go func() {
		for range time.NewTicker(time.Duration(intervalSec) * time.Second).C {
			s.Mem()
		}
	}()
}

// Mem logs current memory statistics.
func (s *StatsService) Mem() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	log.Debug().
		Uint64("alloc", bToMb(m.Alloc)). // currently allocated (in use)
		// Uint64("allocTotal", bToMb(m.TotalAlloc)). // total allocated (cumulative)
		Uint64("sys", bToMb(m.Sys)).             // obtained from OS
		Uint64("heapAlloc", bToMb(m.HeapAlloc)). // heap memory allocated and still in use
		Uint64("heapIdle", bToMb(m.HeapIdle)).   // heap memory not in use
		Uint32("numGC", m.NumGC).                // number of completed GC cycles
		Msg("Memory stats (MB)")
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
