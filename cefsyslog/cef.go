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

package cefsyslog

import (
	"fmt"
	"strings"
)

const (
	cefPrefix     = "CEF"
	cefVersion    = 0
	cefLayout     = "%s:%d|%s|%s|%s|%s|%s|%d"
	deviceVendor  = "Swiss Learning Hub AG"
	deviceProduct = "LMS"
	deviceVersion = "1.0.0"
)

const (
	ExtSourceAddress  = "src"
	ExtSourceUserName = "suser"
	ExtSourceUserID   = "suid"
	ExtReceiptTime    = "rt"
)

// CEF is a single log entry
type CEF struct {
	DeviceVendor  string
	DeviceProduct string
	DeviceVersion string
	EventClassID  string
	Name          string
	Severity      Priority
	Extension     Extensions
}

// Extensions holds additional CEF key-value pairs
type Extensions map[string]string

// NewCEF returns representation
func NewCEF(classID, name string, severity Priority) *CEF {
	return &CEF{
		DeviceVendor:  deviceVendor,
		DeviceProduct: deviceProduct,
		DeviceVersion: deviceVersion,
		EventClassID:  classID,
		Name:          name,
		Severity:      severity,
		Extension:     map[string]string{},
	}
}

// String returns formatted and escaped string representation
func (f *CEF) String() string {
	s := fmt.Sprintf(
		cefLayout,
		cefPrefix,
		cefVersion,
		f.escape(f.DeviceVendor),
		f.escape(f.DeviceProduct),
		f.escape(f.DeviceVersion),
		f.escape(f.EventClassID),
		f.escape(f.Name),
		f.Severity,
	)
	if f.Extension.String() != "" {
		s += "|" + f.Extension.String()
	}
	return s
}

// escape is used for field value escapes
func (f *CEF) escape(s string) string {
	rep := strings.NewReplacer(
		"\\", "\\\\",
		"|", "\\|",
		"\n", "\\n",
	)
	return rep.Replace(s)
}

// String returns formatted and escaped string representation
func (e Extensions) String() string {
	if len(e) == 0 {
		return ""
	}
	var a []string
	for key, value := range e {
		a = append(a, fmt.Sprintf("%s=%s", e.escape(key), e.escape(value)))
	}
	return strings.Join(a, " ")
}

// escape is used for extension key/value escapes
func (e Extensions) escape(s string) string {
	rep := strings.NewReplacer(
		"\\", "\\\\",
		"\n", "\\n",
		"=", "\\=",
	)
	return rep.Replace(s)
}
