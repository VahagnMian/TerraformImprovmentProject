package main

import (
	"fmt"
	logger "github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type DAG struct {
	nodes map[string][]string
}

func NewDAG() *DAG {
	return &DAG{nodes: make(map[string][]string)}
}

func (dag *DAG) AddEdge(from, to string) {
	dag.nodes[from] = append(dag.nodes[from], to)
}

func (dag *DAG) AddNode(node string) {
	if _, exists := dag.nodes[node]; !exists {
		dag.nodes[node] = []string{}
	}
}

func (dag *DAG) Print(prefix string) {
	for from, tos := range dag.nodes {
		if len(tos) == 0 {
			logger.Debug().Msgf("%s -> []", strings.TrimPrefix(from, prefix))
			//fmt.Printf()
		} else {
			//fmt.Printf("%s -> %v\n", from, tos)
			logger.Debug().Msgf("%s -> %v", from, tos)
		}
	}
}

func (dag *DAG) ToDot(prefix string) string {
	dot := "digraph DAG {\n"
	for from, tos := range dag.nodes {
		if len(tos) == 0 {
			dot += fmt.Sprintf("  \"%s\";", strings.TrimPrefix(from, prefix))
		} else {
			for _, to := range tos {
				dot += fmt.Sprintf("  \"%s\" -> \"%s\";", strings.TrimPrefix(from, prefix), strings.TrimPrefix(to, prefix))
			}
		}
	}
	dot += "}\n"
	return dot
}

func (dag *DAG) TopologicalSort() ([]string, error) {
	inDegree := make(map[string]int)
	for node := range dag.nodes {
		if _, exists := inDegree[node]; !exists {
			inDegree[node] = 0
		}
		for _, neighbor := range dag.nodes[node] {
			inDegree[neighbor]++
		}
	}

	queue := []string{}
	for node, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, node)
		}
	}

	sorted := []string{}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		sorted = append(sorted, node)

		for _, neighbor := range dag.nodes[node] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	if len(sorted) != len(dag.nodes) {
		return nil,
			fmt.Errorf("Graph has a cycle")
	}

	return sorted, nil
}

func (dag *DAG) Apply() {
	sorted, err := dag.TopologicalSort()
	if err != nil {
		logger.Error().Msgf("Error: ", err)
		return
	}
	// Reverse the sorted order to process dependencies first
	for i := len(sorted) - 1; i >= 0; i-- {
		logger.Info().Msgf("Processing node: %v", sorted[i])
	}
}

func containsTfFiles(directoryPath string) bool {
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return false
	}
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".tf" {
			return true
		}
	}
	return false
}

func buildDAG(directoryPath string) (*DAG, error) {
	dag := NewDAG()
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".terraform" || info.Name() == ".git" {
				return filepath.SkipDir
			}
			// Check if the directory itself contains .tf files
			if containsTfFiles(path) {
				dag.AddNode(path)
				files, err := ioutil.ReadDir(path)
				if err != nil {
					return err
				}
				for _, file := range files {
					if !file.IsDir() && filepath.Ext(file.Name()) == ".tf" {
						content, err := ioutil.ReadFile(filepath.Join(path, file.Name()))
						if err != nil {
							return err
						}
						referencedDir := getReferencedDirectory(string(content))
						if referencedDir != "" {
							dag.AddEdge(path, filepath.Join(getParentDirectory(path), referencedDir))
						}
					}
				}
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dag, nil
}

func getReferencedDirectory(content string) string {
	re := regexp.MustCompile(`getValueByKey\("([^"]+)", "([^"]+)"\)`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 2 {
		return matches[1]
	}
	return ""
}
