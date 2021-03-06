package rabbitmq

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		DisabledByDefault: true,
		Create:            func() module.Module { return New() },
	}

	module.Register("rabbitmq", creator)
}

const (
	defaultURL         = "http://localhost:15672"
	defaultUsername    = "guest"
	defaultPassword    = "guest"
	defaultHTTPTimeout = time.Second
)

// New creates Rabbitmq with default values
func New() *Rabbitmq {
	return &Rabbitmq{
		HTTP: web.HTTP{
			Request: web.Request{
				URL:      defaultURL,
				Username: defaultUsername,
				Password: defaultPassword,
			},
			Client: web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
	}
}

// Rabbitmq rabbitmq module.
type Rabbitmq struct {
	module.Base

	web.HTTP `yaml:",inline"`

	apiClient *apiClient
}

// Cleanup makes cleanup.
func (Rabbitmq) Cleanup() {}

// Init makes initialization.
func (r *Rabbitmq) Init() bool {
	if r.URL == "" {
		r.Error("URL is not set")
		return false
	}

	client, err := web.NewHTTPClient(r.Client)

	if err != nil {
		r.Error(err)
		return false
	}

	r.apiClient = &apiClient{
		req:        r.Request,
		httpClient: client,
	}

	r.Debugf("using URL %s", r.URL)
	r.Debugf("using timeout: %s", r.Timeout.Duration)

	return true
}

// Check makes check.
func (r *Rabbitmq) Check() bool {
	return len(r.Collect()) > 0
}

// Charts creates Charts.
func (Rabbitmq) Charts() *Charts {
	return charts.Copy()
}

// Collect collects stats.
func (r *Rabbitmq) Collect() map[string]int64 {
	var (
		overview overview
		node     node
		err      error
	)

	if overview, err = r.apiClient.getOverview(); err != nil {
		r.Error(err)
		return nil
	}

	if node, err = r.apiClient.getNodeStats(); err != nil {
		r.Error(err)
		return nil
	}

	return stm.ToMap(overview, node)
}
