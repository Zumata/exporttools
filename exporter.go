package exporttools

import (
	"log"
	"os"
	"os/signal"

	"github.com/prometheus/client_golang/prometheus"
)

type Exporter interface {

	// to be implemented by custom exporter
	Setup() error
	Close() error

	// implemented by BaseExporter, or custom implementation
	Process()

	// satisfy via GenericCollect & GenericDescribe, or custom implementation
	prometheus.Collector
}

func Export(exporter Exporter) error {
	err := exporter.Setup()
	if err != nil {
		return err
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		err := exporter.Close()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	go exporter.Process()

	prometheus.MustRegister(exporter)
	return nil
}
