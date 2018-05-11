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

package publisher

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/cheggaaa/pb.v1"

	"github.com/palantir/distgo/distgo"
)

type FileInfo struct {
	Path      string
	Bytes     []byte
	Checksums Checksums
}

func NewFileInfo(pathToFile string) (FileInfo, error) {
	bytes, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return FileInfo{}, errors.Wrapf(err, "failed to read file %s", pathToFile)
	}
	return NewFileInfoFromBytes(bytes), nil
}

func NewFileInfoFromBytes(bytes []byte) FileInfo {
	sha1Bytes := sha1.Sum(bytes)
	sha256Bytes := sha256.Sum256(bytes)
	md5Bytes := md5.Sum(bytes)

	return FileInfo{
		Path:  "",
		Bytes: bytes,
		Checksums: Checksums{
			SHA1:   hex.EncodeToString(sha1Bytes[:]),
			SHA256: hex.EncodeToString(sha256Bytes[:]),
			MD5:    hex.EncodeToString(md5Bytes[:]),
		},
	}
}

type Checksums struct {
	SHA1   string
	SHA256 string
	MD5    string
}

func (c Checksums) Match(other Checksums) bool {
	nonEmptyEqual := nonEmptyEqual(c.MD5, other.MD5) || nonEmptyEqual(c.SHA1, other.SHA1) || nonEmptyEqual(c.SHA256, other.SHA256)
	// if no non-empty Checksums are equal, Checksums are not equal
	if !nonEmptyEqual {
		return false
	}

	// if non-empty Checksums are different, treat as suspect and return false
	if nonEmptyDiffer(c.MD5, other.MD5) || nonEmptyDiffer(c.SHA1, other.SHA1) || nonEmptyDiffer(c.SHA256, other.SHA256) {
		return false
	}

	// at least one non-empty checksum was equal, and no non-empty Checksums differed
	return true
}

func nonEmptyEqual(s1, s2 string) bool {
	return s1 != "" && s2 != "" && s1 == s2
}

func nonEmptyDiffer(s1, s2 string) bool {
	return s1 != "" && s2 != "" && s1 != s2
}

var (
	GroupIDFlag = distgo.PublisherFlag{
		Name:        "group-id",
		Description: "the Maven group for the product (overrides value specified in publish configuration)",
		Type:        distgo.StringFlag,
	}
	ConnectionInfoURLFlag = distgo.PublisherFlag{
		Name:        "url",
		Description: "URL for publishing (such as https://repository.domain.com)",
		Type:        distgo.StringFlag,
	}
	ConnectionInfoUsernameFlag = distgo.PublisherFlag{
		Name:        "username",
		Description: "username for authentication",
		Type:        distgo.StringFlag,
	}
	ConnectionInfoPasswordFlag = distgo.PublisherFlag{
		Name:        "password",
		Description: "password for authentication",
		Type:        distgo.StringFlag,
	}
)

func BasicConnectionInfoFlags() []distgo.PublisherFlag {
	return []distgo.PublisherFlag{
		ConnectionInfoURLFlag,
		ConnectionInfoUsernameFlag,
		ConnectionInfoPasswordFlag,
	}
}

type BasicConnectionInfo struct {
	URL      string `yaml:"url,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

func (b *BasicConnectionInfo) SetValuesFromFlags(flagVals map[distgo.PublisherFlagName]interface{}) error {
	if err := SetRequiredStringConfigValue(flagVals, ConnectionInfoURLFlag, &b.URL); err != nil {
		return err
	}
	if err := SetConfigValue(flagVals, ConnectionInfoUsernameFlag, &b.Username); err != nil {
		return err
	}
	if err := SetConfigValue(flagVals, ConnectionInfoPasswordFlag, &b.Password); err != nil {
		return err
	}
	return nil
}

func (b *BasicConnectionInfo) UploadDistArtifacts(productTaskOutputInfo distgo.ProductTaskOutputInfo, baseURL string, artifactExists ArtifactExistsFunc, dryRun bool, stdout io.Writer) (artifactPaths []string, uploadedURLs []string, rErr error) {
	for _, currDistID := range productTaskOutputInfo.Product.DistOutputInfos.DistIDs {
		for _, currArtifactPath := range productTaskOutputInfo.ProductDistArtifactPaths()[currDistID] {
			artifactPaths = append(artifactPaths, currArtifactPath)
			var fi FileInfo
			if !dryRun {
				var err error
				fi, err = NewFileInfo(currArtifactPath)
				if err != nil {
					return nil, nil, err
				}
			} else {
				fi = FileInfo{
					Path: currArtifactPath,
				}
			}
			uploadURL, err := b.UploadFile(fi, baseURL, path.Base(currArtifactPath), artifactExists, dryRun, stdout)
			if err != nil {
				return nil, nil, err
			}
			uploadedURLs = append(uploadedURLs, uploadURL)
		}
	}
	return artifactPaths, uploadedURLs, nil
}

func (b *BasicConnectionInfo) UploadFile(fileInfo FileInfo, baseURL, artifactName string, artifactExists ArtifactExistsFunc, dryRun bool, stdout io.Writer) (rURL string, rErr error) {
	rawUploadURL := strings.Join([]string{baseURL, artifactName}, "/")

	filePath := fileInfo.Path
	if filePath != "" {
		if filepath.IsAbs(filePath) {
			if wd, err := os.Getwd(); err == nil {
				if relPath, err := filepath.Rel(wd, filePath); err == nil {
					filePath = relPath
				}
			}
		}
	}
	if !dryRun && artifactExists != nil && artifactExists(artifactName, fileInfo.Checksums, b.Username, b.Password) {
		errMsgParts := []string{"File"}
		if filePath != "" {
			errMsgParts = append(errMsgParts, filePath)
		}
		errMsgParts = append(errMsgParts, fmt.Sprintf("already exists at %s, skipping upload.\n", rawUploadURL))
		fmt.Fprintf(stdout, strings.Join(errMsgParts, " "))
		return rawUploadURL, nil
	}

	uploadURL, err := url.Parse(rawUploadURL)
	if err != nil {
		return rawUploadURL, errors.Wrapf(err, "failed to parse %s as URL", rawUploadURL)
	}

	uploadMsgParts := []string{"Uploading"}
	if filePath != "" {
		uploadMsgParts = append(uploadMsgParts, filePath)
	}
	uploadMsgParts = append(uploadMsgParts, "to", rawUploadURL)
	distgo.PrintlnOrDryRunPrintln(stdout, fmt.Sprintf(strings.Join(uploadMsgParts, " ")), dryRun)

	if !dryRun {
		header := http.Header{}
		addChecksumToHeader(header, "Md5", fileInfo.Checksums.MD5)
		addChecksumToHeader(header, "Sha1", fileInfo.Checksums.SHA1)
		addChecksumToHeader(header, "Sha256", fileInfo.Checksums.SHA256)

		bar := pb.New(len(fileInfo.Bytes)).SetUnits(pb.U_BYTES)
		bar.Output = stdout
		bar.SetMaxWidth(120)
		bar.Start()
		defer bar.Finish()
		reader := bar.NewProxyReader(bytes.NewReader(fileInfo.Bytes))

		req := http.Request{
			Method:        http.MethodPut,
			URL:           uploadURL,
			Header:        header,
			Body:          ioutil.NopCloser(reader),
			ContentLength: int64(len(fileInfo.Bytes)),
		}
		req.SetBasicAuth(b.Username, b.Password)

		resp, err := http.DefaultClient.Do(&req)
		if err != nil {
			errMsgParts := []string{"failed to upload"}
			if filePath != "" {
				errMsgParts = append(errMsgParts, filePath)
			}
			errMsgParts = append(errMsgParts, "to", rawUploadURL)
			return rawUploadURL, errors.Wrapf(err, strings.Join(errMsgParts, " "))
		}
		defer func() {
			if err := resp.Body.Close(); err != nil && rErr == nil {
				rErr = errors.Wrapf(err, "failed to close response body for URL %s", rawUploadURL)
			}
		}()

		if resp.StatusCode >= http.StatusBadRequest {
			msgParts := []string{"uploading"}
			if filePath != "" {
				msgParts = append(msgParts, filePath)
			}
			msgParts = append(msgParts, fmt.Sprintf("to %s resulted in response %q", rawUploadURL, resp.Status))

			msg := fmt.Sprintf(strings.Join(msgParts, " "))
			if body, err := ioutil.ReadAll(resp.Body); err == nil {
				bodyStr := string(body)
				if bodyStr != "" {
					msg += ":\n" + bodyStr
				}
			}
			return rawUploadURL, fmt.Errorf(msg)
		}
	}
	return rawUploadURL, nil
}

// ArtifactExistsFunc returns true if the specified file with the specified checksums already exists in the destination.
type ArtifactExistsFunc func(dstFileName string, checksums Checksums, username, password string) bool

func addChecksumToHeader(header http.Header, checksumName, checksum string) {
	header.Add(fmt.Sprintf("X-Checksum-%s", checksumName), checksum)
}

func MavenProductPath(productTaskOutputInfo distgo.ProductTaskOutputInfo, groupID string) string {
	return path.Join(strings.Replace(groupID, ".", "/", -1), string(productTaskOutputInfo.Product.ID), productTaskOutputInfo.Project.Version)
}

// GetRequiredGroupID returns the value for the GroupID based on the provided inputs. If the provided flagVals map
// contains a non-empty string value for the GroupIDFlag, that value is used. Otherwise, if the PublishOutputInfo for
// the provided ProductTaskOutputInfo is non-nil, its GroupID value is returned. Returns an empty string if no GroupID
// value is specified.
func GetRequiredGroupID(flagVals map[distgo.PublisherFlagName]interface{}, productTaskOutputInfo distgo.ProductTaskOutputInfo) (string, error) {
	if flagVal, ok := flagVals[GroupIDFlag.Name]; ok {
		if groupIDFlagVal := flagVal.(string); groupIDFlagVal != "" {
			return groupIDFlagVal, nil
		}
	}
	if productTaskOutputInfo.Product.PublishOutputInfo != nil && productTaskOutputInfo.Product.PublishOutputInfo.GroupID != "" {
		return productTaskOutputInfo.Product.PublishOutputInfo.GroupID, nil
	}
	return "", PropertyNotSpecifiedError(GroupIDFlag)
}

func PropertyNotSpecifiedError(flag distgo.PublisherFlag) error {
	return errors.Errorf("%s was not specified -- it must be specified in configuration or using a flag", flag.Name)
}

func SetRequiredStringConfigValues(flagVals map[distgo.PublisherFlagName]interface{}, flagAndStringPtrs ...interface{}) error {
	if len(flagAndStringPtrs)%2 != 0 {
		return errors.Errorf("flagsAndStringPtrs parameters must be specified in pairs, got %d", len(flagAndStringPtrs))
	}
	for i := 0; i < len(flagAndStringPtrs); i += 2 {
		if err := SetRequiredStringConfigValue(flagVals, flagAndStringPtrs[i].(distgo.PublisherFlag), flagAndStringPtrs[i+1].(*string)); err != nil {
			return err
		}
	}
	return nil
}

func SetRequiredStringConfigValue(flagVals map[distgo.PublisherFlagName]interface{}, flag distgo.PublisherFlag, stringPtr *string) error {
	if err := SetConfigValue(flagVals, flag, stringPtr); err != nil {
		return err
	}
	if *stringPtr == "" {
		return PropertyNotSpecifiedError(flag)
	}
	return nil
}

func SetConfigValues(flagVals map[distgo.PublisherFlagName]interface{}, flagAndPtrs ...interface{}) error {
	if len(flagAndPtrs)%2 != 0 {
		return errors.Errorf("flagAndPtrs parameters must be specified in pairs, got %d", len(flagAndPtrs))
	}
	for i := 0; i < len(flagAndPtrs); i += 2 {
		if err := SetConfigValue(flagVals, flagAndPtrs[i].(distgo.PublisherFlag), flagAndPtrs[i+1]); err != nil {
			return err
		}
	}
	return nil
}

func SetConfigValue(flagVals map[distgo.PublisherFlagName]interface{}, flag distgo.PublisherFlag, configValPtr interface{}) error {
	configValPtrType := reflect.TypeOf(configValPtr)
	if configValPtrType.Kind() != reflect.Ptr {
		return errors.Errorf("configValPtr type %q is not a pointer type", configValPtrType)
	}

	flagVal, ok := flagVals[flag.Name]
	if !ok {
		// flag is not set: nothing to do
		return nil
	}

	flagValType := reflect.TypeOf(flagVal)
	if elemType := configValPtrType.Elem(); !flagValType.AssignableTo(elemType) {
		return errors.Errorf("flagValType is not assignable to the provided configValPtr: !%v.AssignableTo(%v)", flagValType, elemType)
	}
	// set value pointed at by configValPtr to be the value stored in the provided flags
	reflect.Indirect(reflect.ValueOf(configValPtr)).Set(reflect.ValueOf(flagVal))
	return nil
}
