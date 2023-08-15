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
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/swisslearninghub/logsync/cefsyslog"
	"os"
	"path/filepath"
)

// Config wraps up configured detections
type Config struct {
	Syslog struct {
		Address  string             `json:"address"    validate:"required,hostname_port"`
		Proto    string             `json:"proto"      validate:"required,oneof=tcp udp"`
		Tag      string             `json:"tag"        validate:"required,gt=0,lte=32"`
		Facility cefsyslog.Priority `json:"facility"   validate:"required,oneof=0 8 16 24 32 40 48 56 64 72 80 88"`
	} `json:"syslog"`
	OAuth2 struct {
		ClientID   string `json:"client_id"   validate:"required,gt=0"`
		Secret     string `json:"secret"      validate:"required,gt=0"`
		TokenURL   string `json:"token_url"   validate:"required,url"`
		ContextURL string `json:"context_url" validate:"required,url"`
	} `json:"oauth2"`
	Filter struct {
		Type []string `json:"type"`
		Days int      `json:"days" validate:"required,gt=0,lte=7"`
		Max  int      `json:"max"  validate:"required,gt=0,lte=999999"`
	} `json:"filter"`
	Detections []Detection `json:"detections" validate:"required,gt=0"`
	Logfile    string      `json:"logfile"    validate:"omitempty,gt=0"`
}

var ErrConfigNotFound = errors.New("config not found")

// NewFromFiles returns *Config
func NewFromFiles(paths ...string) (*Config, error) {
	for _, path := range paths {
		if path == "" {
			continue
		}
		cfg, err := NewFromFile(path)
		if err != nil {
			if errors.Is(err, ErrConfigNotFound) {
				continue
			}
			return nil, err
		}
		return cfg, nil
	}
	return nil, ErrConfigNotFound
}

// NewFromFile returns config
func NewFromFile(file string) (*Config, error) {

	var p string
	var err error

	if p, err = filepath.Abs(file); err != nil {
		return nil, err
	}

	if _, err = os.Stat(p); err != nil {
		return nil, ErrConfigNotFound
	}

	var bs []byte

	if bs, err = os.ReadFile(p); err != nil {
		return nil, err
	}

	return NewFromBytes(bs)
}

// NewFromBytes returns Config
func NewFromBytes(bs []byte) (*Config, error) {

	var err error

	c := new(Config)

	if err = json.Unmarshal(bs, c); err != nil {
		return nil, err
	}

	v := validator.New()
	if err = v.Struct(c); err != nil {
		return nil, err
	}

	return c, nil
}
