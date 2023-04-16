package monitor

import (
	"fmt"
	"runtime/metrics"
	"strings"
	"testing"
)

func TestGetMetricsInfo(t *testing.T) {
	var sysinfoMetrics []metrics.Sample
	for _, d := range metrics.All() {
		sysinfoMetrics = append(sysinfoMetrics, metrics.Sample{Name: d.Name})
	}
	if len(sysinfoMetrics) > 0 {
		metrics.Read(sysinfoMetrics)
	}
	//bucketsMap := make(map[string][]float64)
	for _, m := range sysinfoMetrics {
		kind := m.Value.Kind()
		switch kind {
		case metrics.KindFloat64Histogram:
			fmt.Print(m.Name + ": ")
			fmt.Print("KindFloat64Histogram: ")
			//bucketsMap[m.Name] = m.Value.Float64Histogram().Buckets
			unit := m.Name[strings.IndexRune(m.Name, ':')+1:]
			fmt.Println(RuntimeMetricsBucketsForUnit(m.Value.Float64Histogram().Buckets, unit))
		case metrics.KindUint64:
			fmt.Print(m.Name + ": ")
			fmt.Println(m.Value.Uint64())
		case metrics.KindFloat64:
			fmt.Print(m.Name + ": ")
			fmt.Println(m.Value.Float64())
		}
	}

}
