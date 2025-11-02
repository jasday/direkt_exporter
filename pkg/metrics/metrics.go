package metrics

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type Gauge struct {
	Name   string
	Desc   string
	Labels []string
}

func NewGaugeMap(metrics []Gauge) map[string]*prometheus.GaugeVec {
	ret := make(map[string]*prometheus.GaugeVec)
	for _, metric := range metrics {
		ret[metric.Name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metric.Name,
			Help: metric.Desc,
		}, metric.Labels)
	}
	return ret
}

func StringBoolToInt(s string) int {
	if s == "ok" {
		s = "true"
	}
	b, err := strconv.ParseBool(s)
	if b && err == nil {
		return 1
	}
	return 0
}

func BoolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func BoolToFloat64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
