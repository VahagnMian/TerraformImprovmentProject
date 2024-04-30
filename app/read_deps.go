package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadTerraformDependencies(baseDir string) []string {

	var terraformTemplatesPath = []string{}

	// Define which files should be ignored
	excludeFiles := []string{".terraform", ".terraform.lock.hcl", "terraform.tfstate", "terraform.tfstate.backup"}

	// Use filepath.Walk to navigate through the directory tree
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// if this cycle find the file that is excluded, it will be skipped
		for _, v := range excludeFiles {

			if strings.Contains(path, v) {
				return nil
			}

		}

		// Skip all directories
		if info.IsDir() {
			return nil
		}

		// Print the path of the file

		terraformTemplatesPath = append(terraformTemplatesPath, path)
		//fmt.Println(terraformTemplatesPath)
		//fmt.Println(path)
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %v\n", err)
	}

	return terraformTemplatesPath
}
