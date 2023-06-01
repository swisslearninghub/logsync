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

package api

// EventRepresentation is a representation of an event
type EventRepresentation struct {
	Time      int64             `json:"time,omitempty"`
	Type      *string           `json:"type,omitempty"`
	RealmID   *string           `json:"realmId,omitempty"`
	ClientID  *string           `json:"clientId,omitempty"`
	UserID    *string           `json:"userId,omitempty"`
	SessionID *string           `json:"sessionId,omitempty"`
	IPAddress *string           `json:"ipAddress,omitempty"`
	Details   map[string]string `json:"details,omitempty"`
}

// HasDetail returns true if given key exists in details
func (r *EventRepresentation) HasDetail(key string) bool {
	if r.Details == nil {
		return false
	}
	_, ok := r.Details[key]
	return ok
}

// GetDetail returns value if given key exists in details. Returns defaultValue otherwise.
func (r *EventRepresentation) GetDetail(key, defaultValue string) string {
	if !r.HasDetail(key) {
		return defaultValue
	}
	return r.Details[key]
}
