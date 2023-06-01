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

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const httpTimeout = 10 * time.Second

const (
	QueryParamFrom  = "dateFrom"
	QueryParamTo    = "dateTo"
	QueryParamMax   = "max"
	QueryParamType  = "type"
	EventMax        = "999999"
	EventDateLayout = "2006-01-02"
)

// HubAPI is used to handle OAuth2 requests against core-events endpoint
type HubAPI struct {
	client      *http.Client
	tokenSource oauth2.TokenSource
	mu          sync.Mutex
	tok         *oauth2.Token
	contextURL  string
}

func NewAPI(clientID, clientSecret, tokenURL, contextURL string) (*HubAPI, error) {

	a := new(HubAPI)
	a.client = &http.Client{
		Timeout: httpTimeout,
	}
	a.contextURL = contextURL

	ctx := context.TODO()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, a.client)

	conf := &clientcredentials.Config{
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		Scopes:         []string{},
		TokenURL:       tokenURL,
		EndpointParams: url.Values{},
	}

	a.tokenSource = conf.TokenSource(ctx)

	return a, nil
}

// QueryClientEvents builds and executes request from params
func (api *HubAPI) QueryClientEvents(params url.Values) ([]EventRepresentation, error) {

	var baseURL *url.URL
	var req *http.Request
	var token *oauth2.Token
	var err error

	if baseURL, err = url.Parse(api.contextURL + "/client"); err != nil {
		return nil, err
	}

	baseURL.RawQuery = params.Encode()

	if req, err = http.NewRequest(http.MethodGet, baseURL.String(), nil); err != nil {
		return nil, err
	}

	if token, err = api.token(); err != nil {
		return nil, err
	}
	token.SetAuthHeader(req)

	var res *http.Response
	if res, err = api.client.Do(req); err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from api - expected code 200; got code %d", res.StatusCode)
	}

	var bs []byte
	if bs, err = ioutil.ReadAll(res.Body); err != nil {
		return nil, err
	}

	var events []EventRepresentation
	if err = json.Unmarshal(bs, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// token asserts we have a valid token
func (api *HubAPI) token() (*oauth2.Token, error) {

	api.mu.Lock()
	defer api.mu.Unlock()

	if api.tok != nil && api.tok.Valid() {
		return api.tok, nil
	}

	var err error
	if api.tok, err = api.tokenSource.Token(); err != nil {
		return nil, err
	}

	return api.tok, nil
}
