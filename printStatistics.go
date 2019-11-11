package main

import (
	"fmt"
	"logbench/benchmark"
	"runtime"
	"time"

	"github.com/dustin/go-humanize"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func printStatistics(
	logger,
	operation string,
	target uint64,
	concWriters uint,
	stats benchmark.Statistics,
	memStatChan chan MemStats,
) {
	numPrint := message.NewPrinter(language.English)

	memStats := <-memStatChan

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	timeAvg := stats.TotalTime / time.Duration(target)
	totalLogs := numPrint.Sprintf("%d", stats.TotalLogsWritten)
	totalWriters := numPrint.Sprintf("%d", concWriters)
	totalGC := numPrint.Sprintf("%d", m.NumGC)
	totalMallocs := numPrint.Sprintf("%d", m.Mallocs)
	maxHeapObj := numPrint.Sprintf("%d", memStats.MaxHeapObjects)
	heapAllocInc := numPrint.Sprintf("%d", memStats.HeapAllocInc)
	heapObjInc := numPrint.Sprintf("%d", memStats.HeapObjectsInc)
	memStatSamples := numPrint.Sprintf("%d", memStats.StatSamples)

	fmt.Println("")
	fmt.Printf(" logger:        %s\n", logger)
	fmt.Printf(" operation:     %s\n", operation)
	fmt.Printf(" logs total:    %s\n", totalLogs)
	fmt.Printf(" conc. writers: %s\n", totalWriters)
	fmt.Println("")
	fmt.Printf(" time total:    %s\n", stats.TotalTime)
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
