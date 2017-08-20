// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package publish

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type AlmanacUnit struct {
	Product  string            `json:"product"`
	Branch   string            `json:"branch"`
	Revision string            `json:"revision"`
	URL      string            `json:"url"`
	Tags     []string          `json:"tags,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type AlmanacInfo struct {
	URL      string
	AccessID string
	Secret   string
	Release  bool
}

func (a *AlmanacInfo) CheckConnectivity(client *http.Client) error {
	_, err := a.get(client, "/v1/units")
	return err
}

func (a *AlmanacInfo) CheckProduct(client *http.Client, product string) error {
	_, err := a.get(client, strings.Join([]string{"/v1/units", product}, "/"))
	return err
}

func (a *AlmanacInfo) CreateProduct(client *http.Client, product string) error {
	_, err := a.do(client, http.MethodPost, "/v1/units/products", fmt.Sprintf(`{"name":"%s"}`, product))
	return err
}

func (a *AlmanacInfo) CheckProductBranch(client *http.Client, product, branch string) error {
	_, err := a.get(client, strings.Join([]string{"/v1/units", product, branch}, "/"))
	return err
}

func (a *AlmanacInfo) CreateProductBranch(client *http.Client, product, branch string) error {
	_, err := a.do(client, http.MethodPost, strings.Join([]string{"/v1/units", product}, "/"), fmt.Sprintf(`{"name":"%s"}`, branch))
	return err
}

func (a *AlmanacInfo) GetUnit(client *http.Client, product, branch, revision string) ([]byte, error) {
	return a.get(client, strings.Join([]string{"/v1/units", product, branch, revision}, "/"))
}

func (a *AlmanacInfo) CreateUnit(client *http.Client, unit AlmanacUnit, version string, stdout io.Writer) error {
	endpoint := "/v1/units"

	// set version field of metadata to be version
	if unit.Metadata == nil {
		unit.Metadata = make(map[string]string)
	}
	unit.Metadata["version"] = version

	jsonBytes, err := json.Marshal(unit)
	if err != nil {
		return errors.Wrapf(err, "Failed to marshal %v as JSON", unit)
	}

	fmt.Fprintf(stdout, "Creating Almanac unit for product %s, branch %s, revision %s\n", unit.Product, unit.Branch, unit.Revision)
	if _, err := a.do(client, http.MethodPost, endpoint, string(jsonBytes)); err != nil {
		return err
	}
	return nil
}

func (a *AlmanacInfo) ReleaseProduct(client *http.Client, product, branch, revision string) error {
	gaBody := map[string]string{
		"name": "GA",
	}
	jsonBytes, err := json.Marshal(gaBody)
	if err != nil {
		return errors.Wrapf(err, "Failed to marshal %v as JSON", gaBody)
	}

	_, err = a.do(client, http.MethodPost, strings.Join([]string{"/v1/units", product, branch, revision, "releases"}, "/"), string(jsonBytes))
	return err
}

func (a *AlmanacInfo) get(client *http.Client, endpoint string) ([]byte, error) {
	return a.do(client, http.MethodGet, endpoint, "")
}

func (a *AlmanacInfo) do(client *http.Client, method, endpoint, body string) (rBody []byte, rErr error) {
	destURL, err := url.Parse(a.URL + endpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed")
	}

	req := http.Request{
		Method: method,
		URL:    destURL,
	}

	if body != "" {
		req.Header = http.Header{
			"Content-Type": []string{"application/json"},
		}
		req.Body = ioutil.NopCloser(bytes.NewReader([]byte(body)))
		req.ContentLength = int64(len([]byte(body)))
	}

	if err := addAlmanacAuthForRequest(a.AccessID, a.Secret, body, &req); err != nil {
		return nil, errors.Wrapf(err, "Failed to add Almanac authorization info to header for request %v", req)
	}

	resp, err := client.Do(&req)
	if err != nil {
		// remove authorization information from header before returning as part of error
		req.Header.Del("X-authorization")
		return nil, errors.Wrapf(err, "Almanac request failed: %v", req)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil && rErr == nil {
			rErr = errors.Wrapf(err, "failed to close response body for %s", destURL.String())
		}
	}()
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read response body for %s", destURL.String())
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return responseBytes, errors.Errorf("Received non-success status code: %s. Response: %s", resp.Status, string(responseBytes))
	}

	return responseBytes, nil
}

// Adds the X-timestamp and X-authorization header entries to the provided http.Request. Assumes that the body of the
// request will be the byte representation of the "body" string.
func addAlmanacAuthForRequest(accessID, secret, body string, req *http.Request) error {
	// if request Header is nil, initialize to empty object so that map assignment can be done
	if req.Header == nil {
		req.Header = http.Header{}
	}

	timestamp := time.Now().Unix()
	req.Header.Add("X-timestamp", fmt.Sprintf("%d", timestamp))

	hmac, err := hmacSHA1(fmt.Sprint(req.URL.String(), timestamp, body), secret)
	if err != nil {
		return err
	}
	req.Header.Add("X-authorization", fmt.Sprintf("%s:%s", accessID, hmac))
	return nil
}

func hmacSHA1(message string, secret string) (string, error) {
	h := hmac.New(sha1.New, []byte(secret))
	if _, err := h.Write([]byte(message)); err != nil {
		return "", errors.Wrapf(err, "Failed to compute HMAC-SHA1")
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
