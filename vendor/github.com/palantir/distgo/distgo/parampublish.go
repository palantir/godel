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

package distgo

type PublisherTypeID string

type ByPublisherTypeID []PublisherTypeID

func (a ByPublisherTypeID) Len() int           { return len(a) }
func (a ByPublisherTypeID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPublisherTypeID) Less(i, j int) bool { return a[i] < a[j] }

type PublishParam struct {
	// GroupID is the Maven group ID used for the publish operation.
	GroupID string

	// PublishInfo contains extra configuration for the publish operation. The key is the type of publish.
	PublishInfo map[PublisherTypeID]PublisherParam
}

type PublisherParam struct {
	// the raw YAML configuration for this publish operation
	ConfigBytes []byte
}

type PublishOutputInfo struct {
	GroupID string `json:"groupId"`
}

func (p *PublishParam) ToPublishOutputInfo() PublishOutputInfo {
	return PublishOutputInfo{
		GroupID: p.GroupID,
	}
}
