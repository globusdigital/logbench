package main

import (
	"os"
	"runtime"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/globusdigital/logbench/benchmark"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func printStatistics(
	timeTotal time.Duration,
	target uint64,
	concWriters uint,
	memStatChan chan benchmark.MemStats,
	loggerOrder []string,
	operationsOrder []string,
	stats map[string]map[string]benchmark.Statistics,
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
		dr("target", numPrint.Sprintf("%d", target))
		dr("conc. writers", totalConcWriters)
		dr("", "")
		dr("time total", timeTotal.String())
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
			"operation",
			"time total",
			"time avg.",
			"written",
		})
		tbMain.SetAlignment(tablewriter.ALIGN_LEFT)

		for _, loggerName := range loggerOrder {
			for _, operation := range operationsOrder {
				stats := stats[loggerName][operation]
				tbMain.Append([]string{
					loggerName,
					operation,
					stats.TotalTime.String(),
					(stats.TotalTime / time.Duration(target)).String(),
					numPrint.Sprintf("%d", stats.TotalLogsWritten),
				})
			}
		}
		tbMain.Render()
	}
}
