package main

import (
	"fmt"
)

func dfs(graph map[string][]string, start string, path []string, visited map[string]bool, results *[]string) {
	// Check if already visited to handle cycles
	if visited[start] {
		return
	}
	// Mark this node as visited
	visited[start] = true
	// Append this node to the path
	path = append(path, start)

	// Check if this node leads to others
	if len(graph[start]) == 0 {
		// If no children, we are at the end of a path, print it
		*results = append(*results, fmt.Sprintf("%s", path))
	} else {
		// Otherwise, continue to all children
		for _, neighbor := range graph[start] {
			dfs(graph, neighbor, path, visited, results)
		}
	}

	// Backtrack: remove current node from path and mark as not visited
	path = path[:len(path)-1]
	visited[start] = false
}

func findDependencyPaths(graph map[string][]string) []string {
	visited := make(map[string]bool)
	var results []string
	for node := range graph {
		if !visited[node] {
			dfs(graph, node, nil, visited, &results)
		}
	}
	return results
}
