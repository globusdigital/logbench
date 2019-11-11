package main

import (
	"flag"
	"log"
	"logbench/benchmark"
	"logbench/zap"
	"logbench/zerolog"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"sync/atomic"
	"syscall"
	"time"
)

var setups = map[string]benchmark.SetupInit{
	"zap":     zap.Setup,
	"zerolog": zerolog.Setup,
}

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

func setupMemoryWatcher(interval time.Duration) (read chan MemStats) {
	read = make(chan MemStats)
	go func() {
		var stats MemStats
		var runtimeStats runtime.MemStats
		for {
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

			// Wait for the next inspection
			time.Sleep(interval)
		}
	}()
	return
}

func setupTermSigInterceptor() func() bool {
	stop := int32(0)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		for {
			sig := <-sigChan
			if sig == syscall.SIGTERM || sig == syscall.SIGINT {
				// Halt process
				atomic.StoreInt32(&stop, 1)
				break
			}
		}
	}()
	return func() bool { return atomic.LoadInt32(&stop) == 1 }
}

var flagLogger = flag.String("l", "", "logger")
var flagOperation = flag.String("o", benchmark.LogOperationInfo, "operation")
var flagTarget = flag.Uint64(
	"t",
	1_000_000,
	"target number of logs to be written",
)
var flagMemCheckInterval = flag.Duration(
	"mi",
	2*time.Millisecond,
	"memory inspection interval",
)
var flagConcWriters = flag.Uint(
	"w",
	1,
	"number of concurrently writing goroutines",
)

func main() {
	flag.Parse()

	stopped := setupTermSigInterceptor()
	memStatChan := setupMemoryWatcher(*flagMemCheckInterval)

	if *flagLogger == "" {
		log.Fatal("no logger selected")
	}

	// Initialize benchmark function
	setupInit, setupExists := setups[*flagLogger]
	if !setupExists {
		log.Fatalf("no setup for logger %q", *flagLogger)
	}

	bench, err := benchmark.New(os.Stdout, *flagOperation, setupInit)
	if err != nil {
		log.Fatalf("setup %q init: %s", *flagLogger, err)
	}

	stats := bench.Run(*flagTarget, *flagConcWriters, stopped)

	printStatistics(
		*flagLogger,
		*flagOperation,
		*flagTarget,
		*flagConcWriters,
		stats,
		memStatChan,
	)

	// Write memory profile
	runtime.GC()
	memProfile, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer memProfile.Close()
	if err := pprof.WriteHeapProfile(memProfile); err != nil {
		log.Fatal(err)
	}
}
