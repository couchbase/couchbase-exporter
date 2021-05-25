package mocks

import (
	"strings"

	test "github.com/couchbase/couchbase-exporter/test/utils"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
)

type MockSetter struct {
	CallCount     int
	MetricsValues map[string]float64
	MetricsLabels map[string][]string
}

func NewMockSetter() MockSetter {
	setter := MockSetter{
		CallCount:     0,
		MetricsValues: map[string]float64{},
		MetricsLabels: map[string][]string{},
	}

	return setter
}

func (ms *MockSetter) SetGaugeVec(gv prometheus.GaugeVec, value float64, labels ...string) {
	ms.CallCount++

	t, err := gv.GetMetricWithLabelValues(labels...)
	if err != nil {
		return
	}

	var metric io_prometheus_client.Metric

	err = t.Write(&metric)
	if err != nil {
		return
	}

	fqName := test.GetFQNameFromDesc(t.Desc())
	ms.MetricsValues[fqName] = value

	ms.MetricsLabels[fqName] = getLabels(metric)
}

func getLabels(metric io_prometheus_client.Metric) []string {
	simpleLabels := make([]string, 0)
	labels := metric.GetLabel()

	for _, l := range labels {
		value := *l.Value
		simpleLabels = append(simpleLabels, value)
	}

	return simpleLabels
}

func (ms *MockSetter) TestMetric(name string, expectedValue float64, labels ...string) bool {
	v := false
	l := make([]bool, 0)

	for key, value := range ms.MetricsValues {
		if key == name {
			v = value == expectedValue

			for _, label := range ms.MetricsLabels[key] {
				sl := false

				for _, t := range labels {
					if t == label {
						sl = true
						break
					}
				}

				l = append(l, sl)
			}
		}
	}

	return v && reduce(l)
}

func (ms *MockSetter) TestMetricGreaterThanOrEqual(name string, expectedValue float64, labels ...string) bool {
	v := false
	l := make([]bool, 0)

	for key, value := range ms.MetricsValues {
		if strings.HasSuffix(key, name) {
			v = value >= expectedValue

			for _, label := range ms.MetricsLabels[key] {
				sl := false

				for _, t := range labels {
					if t == label {
						sl = true
						break
					}
				}

				l = append(l, sl)
			}
		}
	}

	return v && reduce(l)
}

func reduce(bools []bool) bool {
	result := true

	for _, t := range bools {
		if !t {
			result = false
			break
		}
	}

	return result
}
