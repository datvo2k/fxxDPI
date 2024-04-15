package proxy

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/miekg/dns"
)

type metricValue struct {
	count uint64
}

/*
newMetricValue is a function that creates a new instance of the metricValue struct.

Parameters:
- count (uint64): The initial count value for the metricValue.

Returns:
- *metricValue: A pointer to the newly created metricValue instance.

Example:

	count := uint64(10)
	metric := newMetricValue(count)

Note:

	The newMetricValue function is used to initialize a metricValue struct with an initial count value.
*/
func newMetricValue(count uint64) *metricValue {
	return &metricValue{
		count: count,
	}
}

/*
incrementCount is a method of the metricValue struct.

Parameters:
- None

Returns:
- None

Example:

	metric := &metricValue{}
	metric.incrementCount()

Note:

	The incrementCount method is used to increment the count value of the metricValue instance by 1.
*/
func (metricValue *metricValue) incrementCount() {
	atomic.AddUint64(&(metricValue.count), 1)
}

/*
loadCount is a method of the metricValue struct.

Parameters:
- None

Returns:
- uint64: The current count value of the metricValue instance.

Example:

	metric := &metricValue{}
	count := metric.loadCount()

Note:

	The loadCount method is used to retrieve the current count value of the metricValue instance.
*/
func (metricValue *metricValue) loadCount() uint64 {
	return atomic.LoadUint64(&(metricValue.count))
}

type metrics struct {
	configuration            *MetricsConfiguration
	blockedValue             metricValue
	dohClientErrorsValue     metricValue
	writeResponseErrorsValue metricValue
	rcodeMetricsMap          sync.Map
	rrTypeMetricsMap         sync.Map
}

/*
newMetrics is a function that creates a new instance of the metrics struct.

Parameters:
- configuration (*MetricsConfiguration): The configuration for the metrics.

Returns:
- *metrics: A pointer to the newly created metrics instance.

Example:

	configuration := &MetricsConfiguration{}
	metrics := newMetrics(configuration)

Note:

	The newMetrics function is used to initialize a metrics struct with the given configuration.
*/
func newMetrics(configuration *MetricsConfiguration) *metrics {
	return &metrics{
		configuration: configuration,
	}
}

func (metrics *metrics) incrementBlocked() {
	metrics.blockedValue.incrementCount()
}

func (metrics *metrics) blocked() uint64 {
	return metrics.blockedValue.loadCount()
}

func (metrics *metrics) incrementDOHClientErrors() {
	metrics.dohClientErrorsValue.incrementCount()
}

func (metrics *metrics) dohClientErrors() uint64 {
	return metrics.dohClientErrorsValue.loadCount()
}

func (metrics *metrics) incrementWriteResponseErrors() {
	metrics.writeResponseErrorsValue.incrementCount()
}

func (metrics *metrics) writeResponseErrors() uint64 {
	return metrics.writeResponseErrorsValue.loadCount()
}

func (metrics *metrics) recordRcodeMetric(rcode int) {
	value, loaded := metrics.rcodeMetricsMap.Load(rcode)

	if !loaded {
		value, loaded = metrics.rcodeMetricsMap.LoadOrStore(rcode, newMetricValue(1))
	}

	if loaded {
		value.(*metricValue).incrementCount()
	}
}

/*
rcodeMetricsMapSnapshot is a method of the metrics struct.

Parameters:
- None

Returns:
- map[string]uint64: A map containing the snapshot of the rcode metrics, where the key is the rcode string and the value is the count.

Example:

	metrics := &metrics{}
	snapshot := metrics.rcodeMetricsMapSnapshot()

Note:

	The rcodeMetricsMapSnapshot method is used to retrieve a snapshot of the rcode metrics.
	It iterates over the rcodeMetricsMap and converts the rcode to a string using dns.RcodeToString.
	If the rcode is not found in dns.RcodeToString, it is represented as "UNKNOWN:<rcode>". The count value of each rcode is stored in the map.
*/
func (metrics *metrics) rcodeMetricsMapSnapshot() map[string]uint64 {
	localMap := make(map[string]uint64)
	metrics.rcodeMetricsMap.Range(func(key, value interface{}) bool {
		rcode := key.(int)
		rcodeString, ok := dns.RcodeToString[rcode]
		if !ok {
			rcodeString = fmt.Sprintf("UNKNOWN:%v", rcode)
		}
		rrMetricValue := value.(*metricValue)
		localMap[rcodeString] = rrMetricValue.loadCount()
		return true
	})

	return localMap
}

/*
recordRRTypeMetric is a method of the metrics struct.

Parameters:
- rrType (dns.Type): The DNS resource record type to be recorded.

Returns:
- None

Example:

	metrics := &metrics{}
	metrics.recordRRTypeMetric(dns.TypeA)

Note:

	The recordRRTypeMetric method is used to record the metric for a specific DNS resource record type.
	It checks if the rrType already exists in the rrTypeMetricsMap. If it does, it increments the count value of the metric.
	If it doesn't, it adds the rrType to the rrTypeMetricsMap with an initial count value of 1.
*/
func (metrics *metrics) recordRRTypeMetric(rrType dns.Type) {
	value, loaded := metrics.rrTypeMetricsMap.Load(rrType)

	if !loaded {
		value, loaded = metrics.rrTypeMetricsMap.LoadOrStore(rrType, newMetricValue(1))
	}

	if loaded {
		value.(*metricValue).incrementCount()
	}
}

func (metrics *metrics) rrTypeMetricsMapSnapshot() map[dns.Type]uint64 {

	localMap := make(map[dns.Type]uint64)

	metrics.rrTypeMetricsMap.Range(func(key, value interface{}) bool {
		rrType := key.(dns.Type)
		rrMetricValue := value.(*metricValue)
		localMap[rrType] = rrMetricValue.loadCount()
		return true
	})

	return localMap
}

func (metrics *metrics) String() string {
	return fmt.Sprintf(
		"blocked = %v dohClientErrors = %v writeResponseErrors = %v rcodeMetrics = %v rrtypeMetrics = %v",
		metrics.blocked(), metrics.dohClientErrors(), metrics.writeResponseErrors(),
		metrics.rcodeMetricsMapSnapshot(), metrics.rrTypeMetricsMapSnapshot())
}

func (metrics *metrics) runPeriodicTimer() {
	ticker := time.NewTicker(time.Duration(metrics.configuration.TimerIntervalSeconds) * time.Second)

	for range ticker.C {
		log.Printf("metrics: %v", metrics.String())
	}
}

func (metrics *metrics) start() {
	log.Printf("metrics.start")

	go metrics.runPeriodicTimer()
}
