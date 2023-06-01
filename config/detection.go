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
	"github.com/swisslearninghub/logsync/cefsyslog"
)

// Reporter identifiers
const (
	typeType            = "type"
	typeDetailExists    = "detail_exists"
	typeDetailNotExists = "detail_not_exists"
)

// Detection ...
type Detection struct {
	ClassID   string             `json:"class_id"  validate:"required,gt=0"`
	Name      string             `json:"name"      validate:"required,gt=0"`
	Severity  cefsyslog.Priority `json:"severity"  validate:"required,gte=0,lte=10"`
	LogLevel  cefsyslog.Priority `json:"loglevel"  validate:"required,gte=0,lte=7"`
	Reporters []Reporter         `json:"reporters" validate:"required,gt=0"`
}

// Report returns true if all reporters do (operator: and)
func (d *Detection) Report(er *api.EventRepresentation) bool {

	if len(d.Reporters) == 0 {
		return false
	}

	var rep Report

	reported := 0

	for _, reporter := range d.Reporters {
		switch reporter.Type {
		case typeDetailExists:
			rep = NewDetailExistsReporter(reporter)
		case typeDetailNotExists:
			rep = NewDetailNotExistsReporter(reporter)
		case typeType:
			rep = NewTypeReporter(reporter)
		default:
			continue
		}
		if rep.Do(er) {
			reported++
		}
	}

	return reported == len(d.Reporters)
}

// CEF returns basic *cefsyslog.CEF
func (d *Detection) CEF() *cefsyslog.CEF {
	return cefsyslog.NewCEF(d.ClassID, d.Name, d.Severity)
}
