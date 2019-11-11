package main

import (
	"fmt"
	"logbench/benchmark"
	"os"
	"runtime"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func printStatistics(
	timeTotal time.Duration,
	operation string,
	target uint64,
	concWriters uint,
	memStatChan chan MemStats,
	stats map[string]benchmark.Statistics,
) {
	numPrint := message.NewPrinter(language.English)

	memStats := <-memStatChan

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	totalConcWriters := numPrint.Sprintf("%d", concWriters)
	totalGC := numPrint.Sprintf("%d", m.NumGC)
	totalMallocs := numPrint.Sprintf("%d", m.Mallocs)
	maxHeapObj := numPrint.Sprintf("%d", memStats.MaxHeapObjects)
	heapAllocInc := numPrint.Sprintf("%d", memStats.HeapAllocInc)
	heapObjInc := numPrint.Sprintf("%d", memStats.HeapObjectsInc)
	memStatSamples := numPrint.Sprintf("%d", memStats.StatSamples)

	// Print main table
	{
		tbMain := tablewriter.NewWriter(os.Stdout)
		tbMain.SetAlignment(tablewriter.ALIGN_LEFT)

		dr := func(k, v string) { tbMain.Append([]string{k, v}) }
		dr("operation", operation)
		dr("target", fmt.Sprintf("%d", target))
		dr("", "")
		dr("conc. writers", totalConcWriters)
		dr("max heap", humanize.Bytes(memStats.MaxHeapAlloc))
		dr("max heap obj", maxHeapObj)
		dr("mem samples", memStatSamples)
		dr("heap inc", heapAllocInc)
		dr("heap obj in", heapObjInc)
		dr("heap sys", humanize.Bytes(memStats.MaxHeapSys))
		dr("total alloc", humanize.Bytes(m.TotalAlloc))
		dr("mallocs", totalMallocs)
		dr("num-gc", totalGC)
		dr("total pause", time.Duration(m.PauseTotalNs).String())
		tbMain.Render()
	}

	// Print comparisons table
	{
		tbMain := tablewriter.NewWriter(os.Stdout)
		tbMain.SetHeader([]string{
			"logger",
			"time total",
			"time avg.",
			"written",
		})
		tbMain.SetAlignment(tablewriter.ALIGN_LEFT)

		for loggerName, stats := range stats {
			tbMain.Append([]string{
				loggerName,
				stats.TotalTime.String(),
				(stats.TotalTime / time.Duration(target)).String(),
				numPrint.Sprintf("%d", stats.TotalLogsWritten),
			})
		}
		tbMain.Render()
	}
}
