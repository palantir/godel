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

//// OrderedStringSliceMap represents an ordered map with strings as the keys and a string slice as the values.
//type OrderedStringSliceMap interface {
//	Keys() []string
//	Get(k string) []string
//	Add(k, v string)
//	PutValues(k string, v []string)
//	Remove(k string) bool
//}
//
//type orderedStringSliceMap struct {
//	innerMap    map[string][]string
//	orderedKeys []string
//}
//
//func NewOrderedStringSliceMap() OrderedStringSliceMap {
//	return &orderedStringSliceMap{
//		innerMap: make(map[string][]string),
//	}
//}
//
//func (m *orderedStringSliceMap) Get(k string) []string {
//	return m.innerMap[k]
//}
//
//func (m *orderedStringSliceMap) Add(k, v string) {
//	prevVal, keyPreviouslyInMap := m.innerMap[k]
//	m.innerMap[k] = append(prevVal, v)
//	if !keyPreviouslyInMap {
//		m.orderedKeys = append(m.orderedKeys, k)
//	}
//}
//
//func (m *orderedStringSliceMap) PutValues(k string, v []string) {
//	_, keyPreviouslyInMap := m.innerMap[k]
//	m.innerMap[k] = v
//	if !keyPreviouslyInMap {
//		m.orderedKeys = append(m.orderedKeys, k)
//	}
//}
//
//func (m *orderedStringSliceMap) Remove(k string) bool {
//	if _, keyInMap := m.innerMap[k]; !keyInMap {
//		return false
//	}
//
//	// remove key from map
//	delete(m.innerMap, k)
//	// remove key from slice
//	for i, v := range m.orderedKeys {
//		if v == k {
//			m.orderedKeys = append(m.orderedKeys[:i], m.orderedKeys[i+1:]...)
//			break
//		}
//	}
//	return true
//}
//
//func (m *orderedStringSliceMap) Keys() []string {
//	return m.orderedKeys
//}
//
