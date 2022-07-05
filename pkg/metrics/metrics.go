package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HeadersProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "headers_called",
			Help: "The total number of times headers routed was called",
		},
		[]string{"foo", "foobar"},
	)

	HelloProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hello_called",
		Help: "The total number of times hello route was called",
	})
)
