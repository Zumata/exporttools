package exporttools

import (
	"sort"
	"sync"
)

type FlexMetricStore struct {
	store map[string]*Metric
	mu    *sync.RWMutex
}

func NewFlexMetricStore() *FlexMetricStore {
	return &FlexMetricStore{
		store: make(map[string]*Metric),
		mu:    &sync.RWMutex{},
	}
}

func (fs *FlexMetricStore) MetricNames() []string {
	fs.mu.Lock()
	names := make([]string, 0)
	for k, _ := range fs.store {
		names = append(names, k)
	}
	sort.Strings(names)
	fs.mu.Unlock()
	return names
}

func (fs *FlexMetricStore) Set(m *Metric) error {
	var err error
	fs.mu.Lock()
	if _, ok := fs.store[m.Name]; ok {
		err = fs.store[m.Name].Update(m)
	} else {
		fs.store[m.Name] = m
	}
	fs.mu.Unlock()
	return err
}

func (fs *FlexMetricStore) Get(name string) (*Metric, error) {
	fs.mu.Lock()
	if m, ok := fs.store[name]; ok {
		fs.mu.Unlock()
		return m, nil
	}
	fs.mu.Unlock()
	return &Metric{}, ErrMetricNotFound
}
