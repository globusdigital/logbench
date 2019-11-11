package benchmark

import (
	"context"
	"runtime"
	"time"
)

// MemStats represents memory related statistics
type MemStats struct {
	StatSamples    uint
	HeapAllocInc   uint64
	HeapObjectsInc uint64
	HeapSysInc     uint64
	MaxHeapAlloc   uint64
	MaxHeapObjects uint64
	MaxHeapSys     uint64
}

// StartMemoryWatcher starts a memory watcher goroutine
// which periodically writes memory statistics to the returned channel
func StartMemoryWatcher(
	ctx context.Context,
	interval time.Duration,
) (read chan MemStats) {
	read = make(chan MemStats)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		var stats MemStats
		var runtimeStats runtime.MemStats
		for {
			// Wait for the next inspection
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				stats.StatSamples++

				// Inspect memory usage
				runtime.ReadMemStats(&runtimeStats)

				// Update stats if necessary
				if runtimeStats.HeapAlloc > stats.MaxHeapAlloc {
					stats.HeapAllocInc++
					stats.MaxHeapAlloc = runtimeStats.HeapAlloc
				}
				if runtimeStats.HeapObjects > stats.MaxHeapObjects {
					stats.HeapObjectsInc++
					stats.MaxHeapObjects = runtimeStats.HeapObjects
				}
				if runtimeStats.HeapSys > stats.MaxHeapSys {
					stats.HeapSysInc++
					stats.MaxHeapSys = runtimeStats.HeapSys
				}

				// Try to pass stats to reader if any
				select {
				case read <- stats:
				default:
				}
			}
		}
	}()
	return
}
