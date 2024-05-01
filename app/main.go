package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config []map[string][]interface{}

var workdirPath string = "../workdir/"

func main() {



	files := GetTerraformFiles(workdirPath);


	//fmt.Println(files)

	for _, v := range files {
		dependencies := GetDependency(v)
		if len(dependencies) != 0 {
			for k, v := range dependencies {
				fmt.Print("Path: ", k, " Value: ", v)
			}
		}
	}

	// for _, v := range files {
	// 	fmt.Println(v)
	// }

	
	// refreshTerraformOutputs("workdir/eu-central-1/vpc")
	// //refreshTerraformOutputs("vpc")
	// TerraformTemplateProcessing("workdir/eu-central-1/ec2", "main.tf", false)


	// data, err := ioutil.ReadFile("../vinfra.yaml")
	// if err != nil {
	// 	fmt.Println("Error reading file:", err)
	// 	return
	// }

	// var config Config
	// err = yaml.Unmarshal(data, &config)
	// if err != nil {
	// 	fmt.Println("Error unmarshaling YAML:", err)
	// 	return
	// }

	// Define the root directory where you want to create the structure
	//rootDir := "./workdir"

	// Create and process the directory structure based on the YAML configuration
	// if err := createDirStructure(config, rootDir); err != nil {
	// 	fmt.Println("Error creating directory structure:", err)
	// 	return
	// }

	// fmt.Println("Directory structure created successfully.")
}

func createDirStructure(config Config, basePath string) error {
	for _, item := range config {
		fmt.Println(item)
		for key, val := range item {
			currentPath := filepath.Join(basePath, key)
			if err := os.MkdirAll(currentPath, 0755); err != nil {
				return fmt.Errorf("error creating directory %s: %w", currentPath, err)
			}
			if err := createSubDirs(val, currentPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func createSubDirs(items interface{}, currentPath string) error {
	switch v := items.(type) {
	case []interface{}:
		for _, item := range v {
			if err := processItem(item, currentPath); err != nil {
				return err
			}
		}
	case map[interface{}]interface{}:
		for key, val := range v {
			keyStr, ok := key.(string)
			if !ok {
				return fmt.Errorf("key type not a string, found: %T", key)
			}
			newPath := filepath.Join(currentPath, keyStr)
			if err := os.MkdirAll(newPath, 0755); err != nil {
				return fmt.Errorf("error creating directory %s: %w", newPath, err)
			}
			if err := createSubDirs(val, newPath); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unexpected type: %T", v)
	}
	return nil
}

func processItem(item interface{}, currentPath string) error {
	switch v := item.(type) {
	case string:
		path := filepath.Join(currentPath, v)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %w", path, err)
		}

	case map[string]interface{}:
		for key, val := range v {
			newPath := filepath.Join(currentPath, key)
			if err := os.MkdirAll(newPath, 0755); err != nil {
				return fmt.Errorf("error creating directory %s: %w", newPath, err)
			}
			if err := createSubDirs(val, newPath); err != nil {
				return err
			}
		}
	case map[interface{}]interface{}:
		for key, val := range v {
			keyStr, ok := key.(string)
			if !ok {
				return fmt.Errorf("key type not a string, found: %T", key)
			}
			newPath := filepath.Join(currentPath, keyStr)
			if err := os.MkdirAll(newPath, 0755); err != nil {
				return fmt.Errorf("error creating directory %s: %w", newPath, err)
			}
			if err := createSubDirs(val, newPath); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unhandled item type: %T", v)
	}
	return nil
}
