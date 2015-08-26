package exporttools

import (
	"log"
	"time"
)

type BaseExporter struct {
	Name     string
	Store    MetricStore
	groups   []MetricGroup
	interval time.Duration
}

func NewBaseExporter(name string) *BaseExporter {
	return &BaseExporter{
		Name:     name,
		Store:    NewFlexMetricStore(),
		groups:   []MetricGroup{},
		interval: 15 * time.Second,
	}
}

func (e *BaseExporter) SetInterval(d time.Duration) {
	e.interval = d
}

// AddGroup adds a MetricGroup to the Exporter, for later processing
func (e *BaseExporter) AddGroup(mg MetricGroup) {
	e.groups = append(e.groups, mg)
}

// Process collects all metrics for each MetricGroup
func (e *BaseExporter) Process() {
	for mgIdx := range e.groups {
		go func(mg MetricGroup) {
			ticker := time.NewTicker(e.interval)
			go func() {
				for _ = range ticker.C {
					metrics, err := mg.Collect()
					if err != nil {
						log.Printf("[ERROR] exporter %v during Process. error: %v", e.Name, err)
						continue
					}
					for idx := range metrics {
						e.Store.Set(metrics[idx])
					}
				}
			}()
		}(e.groups[mgIdx])
	}
}
