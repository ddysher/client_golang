// Copyright 2014 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus

// Collector is the interface implemented by anything that can be used by
// Prometheus to collect metrics. The stock metrics provided by this package
// (like Gauge, Counter, Summary) are also MetricCollectors (which only ever
// collect one metric, namely itself). An implementer of Collector may,
// however, collect multiple metrics in a coordinated fashion and/or create
// metrics on the fly. Examples for collectors already implemented in this
// library are the multi-dimensional metrics (i.e. metrics with variable lables)
// like GaugeVec or SummaryVec and the ExpvarCollector.
type Collector interface {
	// Describe returns the super-set of all possible descriptors of
	// metrics collected by this Collector. The returned descriptors
	// fulfill the consistency and uniqueness requirements described in the
	// Desc documentation. This method idempotently returns the same
	// descriptors throughout the lifetime of the Metric.
	Describe() []*Desc
	// Collect is called by Prometheus when collecting metrics. The
	// implementation sends each collected metric via the provided channel
	// and returns once the last metric has been sent. The descriptor of
	// each sent metric is one of those returned by Describe. Returned
	// metrics that share the same descriptor must differ in their variable
	// label values. This method may be called concurrently and must
	// therefore be implemented in a concurrency safe way. Blocking occurs
	// at the expense of total performance of rendering all registered
	// metrics.  Ideally Collector implementations should support concurrent
	// readers.
	Collect(chan<- Metric)
}

// SelfCollector implements Collector for a single metric so that that
// metric collects itself. Add it as an anonymous field to a struct that
// implements Metric, and call Init with the metric itself as an argument.
type SelfCollector struct {
	self  Metric
	descs []*Desc
}

// Init provides the SelfCollector with a reference to the metric it is supposed
// to collect. It is usually called within the factory function to create a
// metric. See example.
func (c *SelfCollector) Init(self Metric) {
	c.self = self
	c.descs = []*Desc{self.Desc()}
}

// Describe implements Collector.
func (c *SelfCollector) Describe() []*Desc {
	return c.descs
}

// Collect implements Collector.
func (c *SelfCollector) Collect(ch chan<- Metric) {
	ch <- c.self
}
