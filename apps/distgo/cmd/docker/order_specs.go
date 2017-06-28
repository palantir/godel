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

package docker

import (
	"sort"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/params"
)

// OrderBuildSpecs orders the provided build specs topologically based on the dependencies among the product specs.
func OrderBuildSpecs(specsWithDeps []params.ProductBuildSpecWithDeps) ([]params.ProductBuildSpecWithDeps, error) {
	var schedule []params.ProductBuildSpecWithDeps
	graph := make(map[string]map[string]struct{})
	specMap := make(map[string][]params.ProductBuildSpecWithDeps)
	// create a graph of dependencies
	for _, curSpec := range specsWithDeps {
		product := curSpec.Spec.ProductName
		specMap[product] = append(specMap[product], curSpec)
		if graph[product] == nil {
			graph[product] = make(map[string]struct{})
		}
		for _, curImage := range curSpec.Spec.DockerImages {
			for depProduct, depTypes := range dockerDepsToMap(curImage.Deps) {
				if !hasDockerDep(depTypes) {
					// only add edge if its a docker image dependency
					continue
				}
				if graph[depProduct] == nil {
					graph[depProduct] = make(map[string]struct{})
				}
				graph[depProduct][product] = struct{}{}
			}
		}
	}

	// get the topological ordering among the products
	order, err := topologicalOrdering(graph)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to generate ordering among the products. The dist dependencies between the products contains a cycle.")
	}

	// construct the final schedule
	for _, product := range order {
		for _, spec := range specMap[product] {
			schedule = append(schedule, spec)
		}
	}

	return schedule, nil
}

func topologicalOrdering(graph map[string]map[string]struct{}) ([]string, error) {
	var order []string
	// get all nodes in the graph and sort lexicographically for deterministic order
	var nodes []string
	indeg := make(map[string]int)
	for node := range graph {
		indeg[node] = 0
		nodes = append(nodes, node)
	}
	sort.Strings(nodes)
	// compute the incoming edges on each vertex
	for _, v := range nodes {
		for neighbor := range graph[v] {
			indeg[neighbor]++
		}
	}
	// q contains all vertices with in-degree zero
	var q []string
	for _, v := range nodes {
		if indeg[v] == 0 {
			q = append(q, v)
		}
	}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		order = append(order, cur)
		var neighbors []string
		// sort all the neighbours to ensure deterministic order
		for neighbor := range graph[cur] {
			neighbors = append(neighbors, neighbor)
		}
		sort.Strings(neighbors)
		for _, neighbor := range neighbors {
			indeg[neighbor]--
			if indeg[neighbor] == 0 {
				q = append(q, neighbor)
			}
		}
	}
	if len(order) != len(graph) {
		return nil, errors.New("Error generating an ordering. Provided DAG contains cyclic dependencies.")
	}
	return order, nil
}

func hasDockerDep(deps map[params.DockerDepType]string) bool {
	for depType := range deps {
		if depType == params.DockerDepDocker {
			return true
		}
	}
	return false
}
