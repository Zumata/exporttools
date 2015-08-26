package exporttools

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

func GenericDescribe(be *BaseExporter, ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc(
		prometheus.BuildFQName("", be.Name, "scrapes"),
		"times exporter has been scraped",
		nil, nil,
	)
	for _, name := range be.Store.MetricNames() {
		m, err := be.Store.Get(name)
		if err != nil {
			log.Printf("[ERROR] exporter %v unable to get metric %v during GenericDescribe due to error: %v", be.Name, name, err)
			continue
		}
		ch <- m.PromDescription(be.Name)
	}
}

func GenericCollect(be *BaseExporter, ch chan<- prometheus.Metric) {
	for _, name := range be.Store.MetricNames() {
		m, err := be.Store.Get(name)
		if err != nil {
			log.Printf("[ERROR] exporter %v unable to get metric %v during GenericCollect due to error: %v", be.Name, name, err)
			continue
		}
		var metric prometheus.Metric
		if m.PromLabels() != nil {
			metric = prometheus.MustNewConstMetric(
				m.PromDescription(be.Name),
				m.PromType(),
				m.PromValue(),
				m.PromLabels()...,
			)
		} else {
			metric = prometheus.MustNewConstMetric(
				m.PromDescription(be.Name),
				m.PromType(),
				m.PromValue(),
			)
		}
		ch <- metric
	}
}
