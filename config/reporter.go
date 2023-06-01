// Copyright 2023 Swiss Learning Hub AG
//
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

package config

import (
	"github.com/swisslearninghub/logsync/api"
	"strings"
)

// Reporter ...
type Reporter struct {
	Type   string            `json:"type"`
	Config map[string]string `json:"config"`
}

// Report to have an entrypoint
type Report interface {
	Do(er *api.EventRepresentation) bool
}

// Has returns true if given configuration key exists
func (r *Reporter) Has(key string) bool {
	if r.Config == nil {
		return false
	}
	_, ok := r.Config[key]
	return ok
}

// Get returns string value of configuration entry. Returns empty string if key does not exist
func (r *Reporter) Get(key string) string {
	if !r.Has(key) {
		return ""
	}
	return r.Config[key]
}

// GetArray returns []string value of configuration entry. Returns nil if key does not exist
func (r *Reporter) GetArray(key, sep string) []string {
	if !r.Has(key) {
		return nil
	}
	if strings.Contains(r.Config[key], sep) {
		return strings.Split(r.Config[key], sep)
	}
	return []string{r.Config[key]}
}

// DetailExistsReporter reports if any of given details are available in event
type DetailExistsReporter struct {
	details []string
}

// NewDetailExistsReporter returns Report
func NewDetailExistsReporter(r Reporter) *DetailExistsReporter {
	return &DetailExistsReporter{details: r.GetArray("details", ",")}
}

// Do match interface
func (ar *DetailExistsReporter) Do(er *api.EventRepresentation) bool {
	for _, detail := range ar.details {
		if er.HasDetail(detail) {
			return true
		}
	}
	return false
}

// DetailNotExistsReporter reports if any of given details are NOT available in event
type DetailNotExistsReporter struct {
	details []string
}

// NewDetailNotExistsReporter returns Report
func NewDetailNotExistsReporter(r Reporter) *DetailNotExistsReporter {
	return &DetailNotExistsReporter{details: r.GetArray("details", ",")}
}

// Do match interface
func (ar *DetailNotExistsReporter) Do(er *api.EventRepresentation) bool {
	for _, detail := range ar.details {
		if !er.HasDetail(detail) {
			return true
		}
	}
	return false
}

// TypeReporter reports if type matches
type TypeReporter struct {
	Type string
}

// NewTypeReporter returns Report
func NewTypeReporter(r Reporter) *TypeReporter {
	return &TypeReporter{Type: r.Get("type")}
}

// Do match interface
func (tr *TypeReporter) Do(er *api.EventRepresentation) bool {
	return *er.Type == tr.Type
}
