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

//go:generate -command generatedocs go run ../../docsgenerator/main.go --base-image "godeltutorial:setup" --tag-prefix godeltutorial --input-dir . --output-dir ../

// Use base image and generate full tutorial:
//go:generate docker build ./baseimage
//go:generate generatedocs

// Use debugimage and generate tutorial starting at step 3:
// go:generate ./debugimage/build.sh
// go:generate generatedocs --start-step=3

package templates
