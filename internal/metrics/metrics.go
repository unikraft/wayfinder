package metrics
// SPDX-License-Identifier: BSD-3-Clause
//
// Authors: Alexander Jung <alex@unikraft.io>
//
// Copyright (c) 2021, Unikraft UG.  All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

import (
  "fmt"
  "time"
	"sync"
	"bytes"
	"encoding/gob"
)

// Measurement represents one measurement value at a timestamp
type Measurement struct {
	Value     []byte
	Timestamp   time.Time
}

// CreateMeasurement creates a new measurement with recent time and bytes from 
// provided value
func CreateMeasurement(value interface{}) (*Measurement, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	if err := enc.Encode(value); err != nil {
		return nil, err
	}

  return &Measurement{
		Value:     buffer.Bytes(),
		Timestamp: time.Now(),
	}, nil
}

// Metric contains a monitoring metric value with current and previous
type Metric struct {
	Name     string
	Values []Measurement
}

// Measurable holds collector metrics in a map
type Measurable struct {
	access  sync.Mutex
	metrics map[string]Metric
}

// NewMeasurable instantiates and returns a new Measurable
func NewMeasurable() *Measurable {
	return &Measurable{
		metrics: make(map[string]Metric),
		access:  sync.Mutex{},
	}
}

// AddMeasurement adds a metric measurement to the series
func (m *Measurable) AddMeasurement(name string, measurement Measurement) {
	m.access.Lock()
	defer m.access.Unlock()

	// create empty metric if not existent
	if _, ok := m.metrics[name]; !ok {
		m.metrics[name] = Metric{
			Name:   name,
			Values: []Measurement{},
		}
	}

	metric := m.metrics[name]

	// add measurement as value to metric
	end := len(metric.Values) - 1
	if end <= 0 {
		end = len(metric.Values)
	}
	metric.Values = append([]Measurement{measurement}, metric.Values[0:end]...)

	// store back
	m.metrics[name] = metric
}

// DeleteMetric removes a metric entirely
func (m *Measurable) DeleteMetric(name string) {
	m.access.Lock()
	defer m.access.Unlock()
	delete(m.metrics, name)
}

// GetMetric reads and returns the metric values by metric name
func (m *Measurable) GetMetric(name string) (*Metric, bool) {
	m.access.Lock()
	defer m.access.Unlock()
	var metric Metric
	metric, exists := m.metrics[name]
	if !exists {
		return &Metric{}, false
	}
	return &metric, true
}

// GetMetricString returns the given metric value at measurement index as string
func (m *Measurable) GetMetricString(name string, index int) string {
	var output string
	if metric, ok := m.GetMetric(name); ok {
		if len(metric.Values) > index {
			decoder := valueDecoder(metric.Values[index].Value)
			var valuetype string
			decoder.Decode(&valuetype)
			output = fmt.Sprintf("%s", valuetype)
		}
	}

	return output
}

// GetMetricIntArray reads and returns a metric int array by metric name
func (m *Measurable) GetMetricIntArray(name string) []int {
	var array []int
	if metric, ok := m.GetMetric(name); ok {
		if len(metric.Values) > 0 {
			decoder := valueDecoder(metric.Values[0].Value)
			decoder.Decode(&array)
		}
	}
	return array
}

// GetMetricStringArray reads and returns a metric string array by metric name
func (m *Measurable) GetMetricStringArray(name string) []string {
	var array []string

	if metric, ok := m.GetMetric(name); ok {
		if len(metric.Values) > 0 {
			decoder := valueDecoder(metric.Values[0].Value)
			decoder.Decode(&array)
		}
	}

	return array
}

// Dump returns all the measurements as map
func (m *Measurable) Dump() map[string]Metric {
	return m.metrics
}

// GetMetricUint64 returns the given metric value at measurement index as string
func (m *Measurable) GetMetricUint64(name string, index int) (string, error) {
	var output string

	if metric, ok := m.GetMetric(name); ok {
		if len(metric.Values) > index {
			decoder := valueDecoder(metric.Values[index].Value)
			var valuetype uint64
			decoder.Decode(&valuetype)
			output = fmt.Sprintf("%d", valuetype)
		} else {
			return "", fmt.Errorf("metric %s has no index %d", name, index)
		}
	} else {
		return "", fmt.Errorf("metric %s not found", name)
	}

  return output, nil
}

// GetMetricFloat64 computes the diff as Float64 between the two measurements
// for given metric and returns it as string
func (m *Measurable) GetMetricFloat64(name string, index int) string {
	var output string

	if metric, ok := m.GetMetric(name); ok {
		if len(metric.Values) > index {
			decoder := valueDecoder(metric.Values[index].Value)
			var valuetype float64
			decoder.Decode(&valuetype)
			output = fmt.Sprintf("%f", valuetype)
		}
	}

  return output
}

// GetMetricDiffUint64AsFloat computes the diff as Uint64 between the two
// measurements for given metric and returns it as string
func (m *Measurable) GetMetricDiffUint64AsFloat(name string, perTime bool) float64 {
	var output float64

	if metric, ok := m.GetMetric(name); ok {
		if len(metric.Values) >= 2 {
			// get first value
			decoder1 := valueDecoder(metric.Values[0].Value)
			var value1 uint64
			decoder1.Decode(&value1)

			// get second value
			decoder2 := valueDecoder(metric.Values[1].Value)
			var value2 uint64
			decoder2.Decode(&value2)

			// calculate value diff per time
			value := float64(value1 - value2)

			// get time diff
			if perTime {
				timeDiff := metric.Values[0].Timestamp.Sub(metric.Values[1].Timestamp).Seconds()
				value = value / timeDiff
			}
			output = value
		}
	}

	return output
}

// GetMetricDiffUint64 computes the diff as Uint64 between the two measurements
// for given metric and returns it as string
func (m *Measurable) GetMetricDiffUint64(name string, perTime bool) string {
	return fmt.Sprintf("%.0f", m.GetMetricDiffUint64AsFloat(name, perTime))
}

func valueDecoder(value []byte) *gob.Decoder {
  return gob.NewDecoder(bytes.NewReader(value))
}
