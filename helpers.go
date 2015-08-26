package exporttools

import (
	"errors"
	"net/http"
)

var (
	ErrIncompatibleMetricType = errors.New("incompatible metric type")
	ErrUnknownMetricType      = errors.New("unknown metric type")
	ErrMetricNotFound         = errors.New("metric does not exist")
)

func DefaultMetricsHandler(title, path string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
      <head><title>` + title + `</title></head>
      <body>
      <h1>` + title + `</h1>
      <p><a href='` + path + `'>Metrics</a></p>
      </body>
      </html>
      `))
	}
}
