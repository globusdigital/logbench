package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/globusdigital/logbench/benchmark"
	"github.com/globusdigital/logbench/logrus"
	"github.com/globusdigital/logbench/zap"
	"github.com/globusdigital/logbench/zerolog"
)

var setups = map[string]benchmark.Setup{
	"zap":     zap.Setup(),
	"zerolog": zerolog.Setup(),
	"logrus":  logrus.Setup(),
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

func main() {
	// Declare and parse flags
	flagLoggers := &flagList{name: "loggers"}
	flag.Var(flagLoggers, "l", "loggers")

	flagOperations := &flagList{name: "operations"}
	flag.Var(flagOperations, "o", "operations")

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
	flagOperationsAll := flag.Bool("o_all", false, "run all operations")

	flag.Parse()

	if *flagOperationsAll {
		flagOperations.vals = []string{
			benchmark.LogOperationInfo,
			benchmark.LogOperationInfoFmt,
			benchmark.LogOperationInfoWithErrorStack,
			benchmark.LogOperationInfoWith3,
			benchmark.LogOperationInfoWith10,
			benchmark.LogOperationInfoWith10Exist,
			benchmark.LogOperationError,
		}
	}

	// Prepare
	stopped := setupTermSigInterceptor()
	memStatChan := benchmark.StartMemoryWatcher(
		context.Background(),
		*flagMemCheckInterval,
	)

	if len(flagLoggers.vals) < 1 {
		log.Fatal("no loggers selected")
	}

	if len(flagOperations.vals) < 1 {
		log.Fatal("no operations selected")
	}

	flagLoggers.RemoveDuplicates()

	stats := make(
		map[string]map[string]benchmark.Statistics,
		len(flagLoggers.vals),
	)

	start := time.Now()
	for _, loggerName := range flagLoggers.vals {
		setupInit, setupExists := setups[loggerName]
		if !setupExists {
			log.Fatalf("no setup for logger %q", loggerName)
		}

		for _, operation := range flagOperations.vals {
			bench, err := benchmark.New(os.Stdout, operation, setupInit)
			if err != nil {
				log.Fatalf("setup %q init: %s", loggerName, err)
			}

			if s := stats[loggerName]; s == nil {
				stats[loggerName] = make(
					map[string]benchmark.Statistics,
					len(flagOperations.vals),
				)
			}
			stats[loggerName][operation] = bench.Run(
				*flagTarget,
				*flagConcWriters,
				stopped,
			)
		}
	}
	timeTotal := time.Since(start)

	printStatistics(
		timeTotal,
		*flagTarget,
		*flagConcWriters,
		memStatChan,
		flagLoggers.vals,
		flagOperations.vals,
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
