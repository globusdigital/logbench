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

var setups = map[string]benchmark.Setup{
	"zap":     zap.Setup(),
	"zerolog": zerolog.Setup(),
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

type flagLoggers []string

func (l *flagLoggers) String() string { return "loggers" }

func (l *flagLoggers) Set(value string) error {
	*l = append(*l, value)
	return nil
}

func (l *flagLoggers) RemoveDuplicates() {
	nw := make([]string, 0, len(*l))
	reg := make(map[string]struct{})
	for _, loggerName := range *l {
		if _, ok := reg[loggerName]; !ok {
			nw = append(nw, loggerName)
			reg[loggerName] = struct{}{}
		}
	}
	*l = nw
}

func main() {
	// Declare and parse flags
	var flagLoggers flagLoggers
	flag.Var(&flagLoggers, "l", "loggers")
	flagOperation := flag.String("o", benchmark.LogOperationInfo, "operation")
	flagTarget := flag.Uint64(
		"t",
		1_000_000,
		"target number of logs to be written",
	)
	flagMemCheckInterval := flag.Duration(
		"mi",
		2*time.Millisecond,
		"memory inspection interval",
	)
	flagConcWriters := flag.Uint(
		"w",
		1,
		"number of concurrently writing goroutines",
	)
	flagMemoryProfile := flag.String(
		"memprof",
		"", // Disabled by default
		"memory profile output file (disabled when empty)",
	)

	flag.Parse()

	// Prepare
	stopped := setupTermSigInterceptor()
	memStatChan := setupMemoryWatcher(*flagMemCheckInterval)

	if len(flagLoggers) < 1 {
		log.Fatal("no loggers selected")
	}

	flagLoggers.RemoveDuplicates()

	stats := make(map[string]benchmark.Statistics, len(flagLoggers))

	start := time.Now()
	for _, loggerName := range flagLoggers {
		// Initialize benchmark function
		setupInit, setupExists := setups[loggerName]
		if !setupExists {
			log.Fatalf("no setup for logger %q", loggerName)
		}

		bench, err := benchmark.New(os.Stdout, *flagOperation, setupInit)
		if err != nil {
			log.Fatalf("setup %q init: %s", loggerName, err)
		}

		stats[loggerName] = bench.Run(*flagTarget, *flagConcWriters, stopped)
	}
	timeTotal := time.Since(start)

	printStatistics(
		timeTotal,
		*flagOperation,
		*flagTarget,
		*flagConcWriters,
		memStatChan,
		flagLoggers,
		stats,
	)

	// Write memory profile
	if *flagMemoryProfile != "" {
		runtime.GC()
		memProfile, err := os.Create(*flagMemoryProfile)
		if err != nil {
			log.Fatal(err)
		}
		defer memProfile.Close()
		if err := pprof.WriteHeapProfile(memProfile); err != nil {
			log.Fatal(err)
		}
	}
}
