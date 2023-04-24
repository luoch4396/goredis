package monitor

import (
	"math"
	"runtime/metrics"
	"strconv"
)

//copy form metrics.All()
const (
	GCHeapTinyAllocsObjects               = "/gc/heap/tiny/allocs:objects"
	GCHeapAllocsObjects                   = "/gc/heap/allocs:objects"
	GCHeapFreesObjects                    = "/gc/heap/frees:objects"
	GCHeapFreesBytes                      = "/gc/heap/frees:bytes"
	GCHeapAllocsBytes                     = "/gc/heap/allocs:bytes"
	GCHeapObjects                         = "/gc/heap/objects:objects"
	GCHeapGoalBytes                       = "/gc/heap/goal:bytes"
	MemoryClassesTotalBytes               = "/memory/classes/total:bytes"
	MemoryClassesHeapObjectsBytes         = "/memory/classes/heap/objects:bytes"
	MemoryClassesHeapUnusedBytes          = "/memory/classes/heap/unused:bytes"
	MemoryClassesHeapReleasedBytes        = "/memory/classes/heap/released:bytes"
	MemoryClassesHeapFreeBytes            = "/memory/classes/heap/free:bytes"
	MemoryClassesHeapStacksBytes          = "/memory/classes/heap/stacks:bytes"
	MemoryClassesOSStacksBytes            = "/memory/classes/os-stacks:bytes"
	MemoryClassesMetadataMSpanInuseBytes  = "/memory/classes/metadata/mspan/inuse:bytes"
	MemoryClassesMetadataMSPanFreeBytes   = "/memory/classes/metadata/mspan/free:bytes"
	MemoryClassesMetadataMCacheInuseBytes = "/memory/classes/metadata/mcache/inuse:bytes"
	MemoryClassesMetadataMCacheFreeBytes  = "/memory/classes/metadata/mcache/free:bytes"
	MemoryClassesProfilingBucketsBytes    = "/memory/classes/profiling/buckets:bytes"
	MemoryClassesMetadataOtherBytes       = "/memory/classes/metadata/other:bytes"
	MemoryClassesOtherBytes               = "/memory/classes/other:bytes"
	Goroutines                            = "/sched/goroutines:goroutines"
)

// GetMetricsInfo 获取Metrics列表
func GetMetricsInfo(metricsNames ...string) map[string]string {
	if len(metricsNames) == 0 {
		return nil
	}
	metricsMap := make(map[string]string, 8)
	var sysInfoMetrics []metrics.Sample
	for _, m := range metricsNames {
		sysInfoMetrics = append(sysInfoMetrics, metrics.Sample{Name: m})
	}
	if len(sysInfoMetrics) > 0 {
		metrics.Read(sysInfoMetrics)
	}

	for _, m := range sysInfoMetrics {
		kind := m.Value.Kind()
		switch kind {
		case metrics.KindBad:
			panic("unexpected or unsupported metric")
		case metrics.KindFloat64Histogram:
			//TODO: 未解析，暂时不需要
			//unit := m.Name[strings.IndexRune(m.Name, ':')+1:]
			//RuntimeMetricsBucketsForUnit(m.Value.Float64Histogram().Buckets, unit)
		case metrics.KindUint64:
			metricsMap[m.Name] = strconv.FormatUint(m.Value.Uint64(), 10)
		case metrics.KindFloat64:
			metricsMap[m.Name] = strconv.FormatFloat(m.Value.Float64(), 'f', 10, 64)
		}
	}
	return metricsMap
}

// RuntimeMetricsBucketsForUnit read and format from buckets
func RuntimeMetricsBucketsForUnit(buckets []float64, unit string) []float64 {
	switch unit {
	case "bytes":
		// Re-bucket as powers of 2.
		return reBucketExp(buckets, 2)
	case "seconds":
		// Re-bucket as powers of 10 and then merge all buckets greater
		// than 1 second into the +Inf bucket.
		b := reBucketExp(buckets, 10)
		for i := range b {
			if b[i] <= 1 {
				continue
			}
			b[i] = math.Inf(1)
			b = b[:i+1]
			break
		}
		return b
	}
	return buckets
}

func reBucketExp(buckets []float64, base float64) []float64 {
	bucket := buckets[0]
	var newBuckets []float64
	// We may see a -Inf here, in which case, add it and skip it
	// since we risk producing NaNs otherwise.
	//
	// We need to preserve -Inf values to maintain runtime/metrics
	// conventions. We'll strip it out later.
	if bucket == math.Inf(-1) {
		newBuckets = append(newBuckets, bucket)
		buckets = buckets[1:]
		bucket = buckets[0]
	}
	// From now on, bucket should always have a non-Inf value because
	// Infs are only ever at the ends of the bucket lists, so
	// arithmetic operations on it are non-NaN.
	for i := 1; i < len(buckets); i++ {
		if bucket >= 0 && buckets[i] < bucket*base {
			// The next bucket we want to include is at least bucket*base.
			continue
		} else if bucket < 0 && buckets[i] < bucket/base {
			// In this case the bucket we're targeting is negative, and since
			// we're ascending through buckets here, we need to divide to get
			// closer to zero exponentially.
			continue
		}
		// The +Inf bucket will always be the last one, and we'll always
		// end up including it here because bucket
		newBuckets = append(newBuckets, bucket)
		bucket = buckets[i]
	}
	return append(newBuckets, bucket)
}
