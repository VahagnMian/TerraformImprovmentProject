package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl"
	logger "github.com/rs/zerolog/log"
)

// getAllOutputs is used to get all outputs from particular module in HCL format
func getAllOutputs(modulePath string, json bool) []byte {

	var cmd *exec.Cmd

	if json {
		cmd = exec.Command("terraform", "output", "-json")
	} else {
		cmd = exec.Command("terraform", "output")
	}

	cmd.Dir = modulePath

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing Terraform:", err)
	}

	return output
}

// parseHCL by passing terraform outputs in HCL format to this function,
// it returns go map of key and value
func parseHCL(hclData string) map[string]interface{} {

	// Create a variable to hold the parsed data
	var result map[string]interface{}

	// Parse the HCL data
	err := hcl.Decode(&result, hclData)
	if err != nil {
		log.Fatalf("Failed to decode HCL: %s", err)
	}

	return result

}

// getValueByKey by passing key and map we can extract exact Value from terraform output
func getValueByKey(key string, result map[string]interface{}) string {

	value := fmt.Sprintf("%v", result[key])

	if checkType(value) == "list" {

	}

	switch checkType(value) {

	case "list":
		value = makeHCLArrayFromHCLOutput(value)
	case "string":
		value = makeHCLStringFromHCLOutput(value)
	}

	return value

}
func refreshTerraformOutputs(modulePath string) {

	var stderr bytes.Buffer
	logger.Info().Msgf("Starting Terraform outputs syncing")

	cmd := exec.Command("terraform", "-chdir=../"+modulePath, "apply", "-refresh-only", "-auto-approve")

	cmd.Stderr = &stderr

	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(string(stdout))
		logger.Error().Msgf("Error occurred during terraform outputs syncing (refresh apply) %v", err)
		fmt.Println(stderr.String())
		return
	}

	logger.Info().Msgf("Terraform outputs refreshed (synced) successfully for " + modulePath + " module")

}

func makeHCLArrayFromHCLOutput(s string) string {

	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")

	elements := strings.Fields(s)

	var hclArray string

	for i := 0; i < len(elements); i++ {
		hclArray = hclArray + "\"" + elements[i] + "\"" + ", "
	}
	hclArray = strings.TrimSuffix(hclArray, ", ")
	hclArray = "[" + hclArray + "]"

	return hclArray

}

func makeHCLStringFromHCLOutput(s string) string {

	value := "\"" + s + "\""

	return value

}

func checkType(s string) string {

	isArray := strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]")

	if isArray {
		return "list"
	}
	return "string"
}

func appendProcessedToTf(fileName string) string {

	result := strings.TrimSuffix(fileName, ".tf")
	result = result + "_processed"
	result = result + ".tf"

	return result

}

func trimProcessedFromTf(fileName string) string {

	result := strings.TrimSuffix(fileName, "_processed.tf")

	result = result + ".tf"

	return result

}


// Get All Terraform files in directory
func GetTerraformFiles(baseDir string) []string {

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


// Get dependencies in terraform files

func GetDependency(tfFile string) map[string]string {

	destinationTemplateFile := tfFile

	// Read file
	file, err := os.Open(tfFile)
	defer file.Close()

	// Error handling
	if err != nil {
		log.Fatal("[ Error ] error occured during reading of: ", destinationTemplateFile, " please check if file exists, and has right permissions")
		os.Exit(1)
	}

	pattern := regexp.MustCompile(`getValueByKey\("([^"]+)", "([^"]+)"\)`)
	scanner := bufio.NewScanner(file)

	lines := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		matches := pattern.FindStringSubmatch(line)

		if matches != nil {

			lines[file.Name()] = line

		}

	}

	return lines

} 