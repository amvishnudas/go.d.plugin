package springboot2

import (
	"github.com/netdata/go-orchestrator/module"
)

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	{
		ID:    "response_codes",
		Title: "Response Codes", Units: "requests/s", Fam: "response_code", Type: module.Stacked,
		Dims: Dims{
			{ID: "resp_1xx", Name: "1xx", Algo: module.Incremental},
			{ID: "resp_2xx", Name: "2xx", Algo: module.Incremental},
			{ID: "resp_3xx", Name: "3xx", Algo: module.Incremental},
			{ID: "resp_4xx", Name: "4xx", Algo: module.Incremental},
			{ID: "resp_5xx", Name: "5xx", Algo: module.Incremental},
		},
	},
	{
		ID:    "thread",
		Title: "Threads", Units: "threads", Fam: "threads", Type: module.Area,
		Dims: Dims{
			{ID: "threads_daemon", Name: "daemon"},
			{ID: "threads", Name: "total"},
		},
	},
	{
		ID:    "heap",
		Title: "Overview", Units: "B", Fam: "heap", Type: module.Stacked,
		Dims: Dims{
			{ID: "mem_free", Name: "free"},
			{ID: "heap_used_eden", Name: "eden"},
			{ID: "heap_used_survivor", Name: "survivor"},
			{ID: "heap_used_old", Name: "old"},
		},
	},
	{
		ID:    "heap_eden",
		Title: "Eden Space", Units: "B", Fam: "heap", Type: module.Area,
		Dims: Dims{
			{ID: "heap_used_eden", Name: "used"},
			{ID: "heap_committed_eden", Name: "committed"},
		},
	},
	{
		ID:    "heap_survivor",
		Title: "Survivor Space", Units: "B", Fam: "heap", Type: module.Area,
		Dims: Dims{
			{ID: "heap_used_survivor", Name: "used"},
			{ID: "heap_committed_survivor", Name: "committed"},
		},
	},
	{
		ID:    "heap_old",
		Title: "Old Space", Units: "B", Fam: "heap", Type: module.Area,
		Dims: Dims{
			{ID: "heap_used_old", Name: "used"},
			{ID: "heap_committed_old", Name: "committed"},
		},
	},
	{
		ID:    "uptime",
		Title: "The uptime of the Java virtual machine", Units: "seconds", Fam: "uptime", Type: module.Line,
		Dims: Dims{
			{ID: "uptime", Name: "uptime", Div: 1000},
		},
	},
}
