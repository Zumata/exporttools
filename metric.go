package exporttools

import "github.com/prometheus/client_golang/prometheus"

type MetricGroup interface {
	Collect() ([]*Metric, error)
}

type MetricStore interface {
	MetricNames() []string
	Set(m *Metric) error
	Get(name string) (*Metric, error)
}

type Metric struct {
	Name        string
	Description string
	Type        metricType
	Value       int64
	LabelKeys   []string
	LabelVals   []string
}

type metricType int

const (
	Counter metricType = iota
	Gauge
)

func (m *Metric) Update(newMetric *Metric) error {
	switch newMetric.Type {
	case Gauge:
		if m.Type != Gauge {
			return ErrIncompatibleMetricType
		}
		m.Value = newMetric.Value
	case Counter:
		if m.Type != Counter {
			return ErrIncompatibleMetricType
		}
		m.Value = m.Value + newMetric.Value
	default:
		return ErrUnknownMetricType
	}
	return nil
}

func (m *Metric) PromDescription(exporterName string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName("", exporterName, m.Name),
		m.Description,
		m.LabelKeys, nil,
	)
}

func (m *Metric) PromType() prometheus.ValueType {
	if m.Type == Counter {
		return prometheus.CounterValue
	}
	return prometheus.GaugeValue
}

func (m *Metric) PromValue() float64 {
	return float64(m.Value)
}

func (m *Metric) PromLabels() []string {
	return m.LabelVals
}
