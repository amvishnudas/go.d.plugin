package springboot2

import (
	"strings"
	"time"

	mtx "github.com/netdata/go.d.plugin/pkg/metrics"
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	module.Register("springboot2", module.Creator{
		Create: func() module.Module { return New() },
	})
}

const (
	defaultHTTPTimeout = time.Second
)

// New returns SpringBoot2 instance with default values
func New() *SpringBoot2 {
	return &SpringBoot2{
		HTTP: web.HTTP{
			Client: web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
	}
}

// SpringBoot2 Spring boot 2 module
type SpringBoot2 struct {
	module.Base

	web.HTTP `yaml:",inline"`

	prom prometheus.Prometheus
}

type metrics struct {
	Uptime mtx.Gauge `stm:"uptime,1000"`

	ThreadsDaemon mtx.Gauge `stm:"threads_daemon"`
	Threads       mtx.Gauge `stm:"threads"`

	Resp1xx mtx.Counter `stm:"resp_1xx"`
	Resp2xx mtx.Counter `stm:"resp_2xx"`
	Resp3xx mtx.Counter `stm:"resp_3xx"`
	Resp4xx mtx.Counter `stm:"resp_4xx"`
	Resp5xx mtx.Counter `stm:"resp_5xx"`

	HeapUsed      heap `stm:"heap_used"`
	HeapCommitted heap `stm:"heap_committed"`

	MemFree mtx.Gauge `stm:"mem_free"`
}

type heap struct {
	Eden     mtx.Gauge `stm:"eden"`
	Survivor mtx.Gauge `stm:"survivor"`
	Old      mtx.Gauge `stm:"old"`
}

// Cleanup Cleanup
func (SpringBoot2) Cleanup() {}

// Init makes initialization
func (s *SpringBoot2) Init() bool {
	client, err := web.NewHTTPClient(s.Client)
	if err != nil {
		s.Error(err)
		return false
	}
	s.prom = prometheus.New(client, s.Request)
	return true
}

// Check makes check
func (s *SpringBoot2) Check() bool {
	rawMetrics, err := s.prom.Scrape()
	if err != nil {
		s.Error(err)
		return false
	}
	jvmMemory := rawMetrics.FindByName("jvm_memory_used_bytes")

	return len(jvmMemory) > 0
}

// Charts creates Charts
func (SpringBoot2) Charts() *Charts {
	return charts.Copy()
}

// Collect collects metrics
func (s *SpringBoot2) Collect() map[string]int64 {
	rawMetrics, err := s.prom.Scrape()
	if err != nil {
		return nil
	}

	var m metrics

	// uptime
	m.Uptime.Set(rawMetrics.FindByName("process_uptime_seconds").Max())

	// response
	gatherResponse(rawMetrics, &m)

	// threads
	m.ThreadsDaemon.Set(rawMetrics.FindByNames("jvm_threads_daemon", "jvm_threads_daemon_threads").Max())
	m.Threads.Set(rawMetrics.FindByNames("jvm_threads_live", "jvm_threads_live_threads").Max())

	// heap memory
	gatherHeap(rawMetrics.FindByName("jvm_memory_used_bytes"), &m.HeapUsed)
	gatherHeap(rawMetrics.FindByName("jvm_memory_committed_bytes"), &m.HeapCommitted)
	m.MemFree.Set(m.HeapCommitted.Sum() - m.HeapUsed.Sum())

	return stm.ToMap(m)
}

func gatherHeap(rawMetrics prometheus.Metrics, m *heap) {
	for _, metric := range rawMetrics {
		id := metric.Labels.Get("id")
		value := metric.Value
		switch {
		case strings.Contains(id, "Eden"):
			m.Eden.Set(value)
		case strings.Contains(id, "Survivor"):
			m.Survivor.Set(value)
		case strings.Contains(id, "Old") || strings.Contains(id, "Tenured"):
			m.Old.Set(value)
		}
	}
}

func gatherResponse(rawMetrics prometheus.Metrics, m *metrics) {
	for _, metric := range rawMetrics.FindByName("http_server_requests_seconds_count") {
		status := metric.Labels.Get("status")
		if status == "" {
			continue
		}
		value := metric.Value
		switch status[0] {
		case '1':
			m.Resp1xx.Add(value)
		case '2':
			m.Resp2xx.Add(value)
		case '3':
			m.Resp3xx.Add(value)
		case '4':
			m.Resp4xx.Add(value)
		case '5':
			m.Resp5xx.Add(value)
		}
	}
}

func (h heap) Sum() float64 {
	return h.Eden.Value() + h.Survivor.Value() + h.Old.Value()
}
