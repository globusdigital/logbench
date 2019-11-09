package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/dustin/go-humanize"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var setups = map[string]func() (FuncSet, error){
	"zerolog": initZerolog,
	"zap":     initZap,
}

const (
	logOperationInfo               = "info"
	logOperationInfoFmt            = "info_fmt"
	logOperationInfoWithErrorStack = "info_with_error_stack"
	logOperationInfoWith3          = "info_with_3"
	logOperationInfoWith10         = "info_with_10"
	logOperationError              = "error"
)

// TimeFormat defines the time logging format
const TimeFormat = "2006-01-02T15:04:05.999999999-07:00"

// Fields3 is a list of 3 fields
type Fields3 struct {
	Name1 string
	Name2 string
	Name3 string

	Value1 string
	Value2 int
	Value3 float64
}

// Fields10 is a list of 10 fields
type Fields10 struct {
	Name1  string
	Name2  string
	Name3  string
	Name4  string
	Name5  string
	Name6  string
	Name7  string
	Name8  string
	Name9  string
	Name10 string

	Value1  string
	Value2  string
	Value3  string
	Value4  string
	Value5  string
	Value6  int
	Value7  float64
	Value8  []string
	Value9  []int
	Value10 []float64
}

// FuncSet defines the callback functions for all benchmarked cases
type FuncSet struct {
	Info               func(msg string)
	InfoFmt            func(msg string, data int)
	Error              func(msg string)
	InfoWithErrorStack func(msg string, err error)
	InfoWith3          func(msg string, fields *Fields3)
	InfoWith10         func(msg string, fields *Fields10)
}

func initSetup(
	setup func() (FuncSet, error),
	operation string,
) (func(), error) {
	FuncSet, err := setup()
	if err != nil {
		return nil, err
	}

	switch operation {
	case logOperationInfo:
		return func() { FuncSet.Info("information") }, nil
	case logOperationInfoFmt:
		return func() { FuncSet.InfoFmt("information", 42) }, nil
	case logOperationInfoWithErrorStack:
		err := errors.New("error with stack trace")
		return func() { FuncSet.InfoWithErrorStack("information", err) }, nil
	case logOperationError:
		return func() { FuncSet.Error("error message") }, nil
	case logOperationInfoWith3:
		fields := &Fields3{
			Name1: "field1", Value1: "some textual value",
			Name2: "field_2_int", Value2: 42,
			Name3: "field_3_float_64", Value3: 42.5,
		}
		return func() { FuncSet.InfoWith3("information", fields) }, nil
	case logOperationInfoWith10:
		fields := &Fields10{
			Name1: "field1", Value1: "",
			Name2: "field2", Value2: "some textual value",
			Name3: "field3", Value3: "and another textual value",
			Name5: "field5", Value5: "an even longer textual value",
			Name6: "field_6_int", Value6: 42,
			Name7: "field_7_float_64", Value7: 42.5,
			Name8: "field_8_multipleStrings", Value8: []string{
				"first",
				"second",
				"third",
			},
			Name9: "field_9_multipleIntegers", Value9: []int{
				0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
			},
			Name10: "field_9_multipleIntegers", Value10: []float64{
				11.5, 24.9, 99.99, 50.5001, 1000.11,
			},
		}
		return func() { FuncSet.InfoWith10("information", fields) }, nil
	}
	return nil, fmt.Errorf("unsupported operation: %q", operation)
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

var flagLogger = flag.String("l", "", "logger")
var flagOperation = flag.String("o", logOperationInfo, "operation")
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

func printStatistics(
	totalLogsWritten uint64,
	timeTotal time.Duration,
	memStatChan chan MemStats,
) {
	totalLogsWritten -= uint64(*flagConcWriters)

	numPrint := message.NewPrinter(language.English)

	memStats := <-memStatChan

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	timeAvg := timeTotal / time.Duration(*flagTarget)
	totalLogs := numPrint.Sprintf("%d", totalLogsWritten)
	totalWriters := numPrint.Sprintf("%d", *flagConcWriters)
	totalGC := numPrint.Sprintf("%d", m.NumGC)
	totalMallocs := numPrint.Sprintf("%d", m.Mallocs)
	maxHeapObj := numPrint.Sprintf("%d", memStats.MaxHeapObjects)
	heapAllocInc := numPrint.Sprintf("%d", memStats.HeapAllocInc)
	heapObjInc := numPrint.Sprintf("%d", memStats.HeapObjectsInc)
	memStatSamples := numPrint.Sprintf("%d", memStats.StatSamples)

	fmt.Println("")
	fmt.Printf(" logger:        %s\n", *flagLogger)
	fmt.Printf(" operation:     %s\n", *flagOperation)
	fmt.Printf(" logs total:    %s\n", totalLogs)
	fmt.Printf(" conc. writers: %s\n", totalWriters)
	fmt.Println("")
	fmt.Printf(" time total:    %s\n", timeTotal)
	fmt.Printf(" time avg:      %s\n", timeAvg)
	fmt.Println("")
	fmt.Printf(" max heap:      %s\n", humanize.Bytes(memStats.MaxHeapAlloc))
	fmt.Printf(" max heap obj:  %s\n", maxHeapObj)
	fmt.Printf(" mem samples:   %s\n", memStatSamples)
	fmt.Printf(" heap inc:      %s\n", heapAllocInc)
	fmt.Printf(" heap obj inc:  %s\n", heapObjInc)
	fmt.Printf(" heap sys:      %s\n", humanize.Bytes(memStats.MaxHeapSys))
	fmt.Printf(" total alloc:   %s\n", humanize.Bytes(m.TotalAlloc))
	fmt.Printf(" mallocs:       %s\n", totalMallocs)
	fmt.Printf(" num-gc:        %s\n", totalGC)
	fmt.Printf(" total pause:   %s\n", time.Duration(m.PauseTotalNs))
	fmt.Println("")
}

func main() {
	flag.Parse()

	stopped := setupTermSigInterceptor()
	memStatChan := setupMemoryWatcher(*flagMemCheckInterval)

	if *flagLogger == "" {
		log.Fatal("no logger selected")
	}

	// Initialize benchmark function
	setupFn, setupExists := setups[*flagLogger]
	if !setupExists {
		log.Fatalf("no setup for logger %q", *flagLogger)
	}

	fn, err := initSetup(setupFn, *flagOperation)
	if err != nil {
		log.Fatalf("setup %q init: %s", *flagLogger, err)
	}

	// Execute benchmark
	start := time.Now()
	logsWritten := uint64(0)

	wg := sync.WaitGroup{}
	wg.Add(int(*flagConcWriters))
	for wk := uint(0); wk < *flagConcWriters; wk++ {
		go func() {
			defer wg.Done()
			for {
				if stopped() {
					break
				}
				if atomic.AddUint64(&logsWritten, 1) > *flagTarget {
					break
				}
				// Write a log
				fn()
			}
		}()
	}
	wg.Wait()

	timeTotal := time.Since(start)

	printStatistics(atomic.LoadUint64(&logsWritten), timeTotal, memStatChan)

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
