## exporttools (ALPHA)
Building blocks for quickly creating custom prometheus exporters.

**exporttools** aims to make it very easy to implement and main exporters via abstracting common patterns of metric collection. These building blocks have been heavily influenced by custom exporter implementations by prometheus.io and other organisations.

#### Getting Started

1. Bootstrap your Exporter by creating a struct and embedding `*BaseExporter`. Your Exporter must satisfy the `Exporter` interface.
```
type Exporter interface {

	// to be implemented by custom exporter
	Setup() error
	Close() error

	// implemented by BaseExporter, or custom implementation
  	// implemented by BaseExporter, or custom implementation
	AddGroup(mg MetricGroup)
	Process()

	// satisfy via GenericCollect & GenericDescribe, or custom implementation
	prometheus.Collector
}
```

2. Implement methods the `Setup()` and `Close()` methods required by the Exporter interface to create/destroy infrastructure connections.

3. For each group of metrics to be collected by the Exporter, satisfy the `MetricGroup` interface then add to your Exporter via `AddGroup()` provided by the `BaseExporter`, within your Exporter's `Setup()` method.
```
type MetricGroup interface {
	Collect() ([]*Metric, error)
}
```

3. `Metrics` handle both counters and gauges
```
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
```

4. Implement `prometheus.Collector` using the helpers:
```
func (e *customExporter) Describe(ch chan<- *prometheus.Desc) {
	exporttools.GenericDescribe(e.BaseExporter, ch)
}

func (e *customExporter) Collect(ch chan<- prometheus.Metric) {
	exporttools.GenericCollect(e.BaseExporter, ch)
}
```

5. Calling `Export(exporter)` will enable collection for all metric groups.
```
func main()
  exporter := postgres.NewCustomExporter()
  err := exporttools.Export(exporter)
  if err != nil {
    log.Fatal(err)
  }
}
```

#### License
MIT

Copyright (c) 2015 Zumata Technologies Pte Ltd.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
