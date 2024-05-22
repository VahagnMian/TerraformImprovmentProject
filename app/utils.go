package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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
		logger.Error().Msgf("Error executing Terraform:", err)
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
func refreshTerraformOutputs(modulePath string) error {

	var stderr bytes.Buffer
	logger.Info().Msgf("Starting Terraform outputs syncing")

	//cmd := exec.Command("terraform", "-chdir=../"+modulePath, "apply", "-refresh-only", "-auto-approve")
	cmd := exec.Command("terraform", "-chdir="+modulePath, "apply", "-auto-approve")

	cmd.Stderr = &stderr

	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(string(stdout))
		logger.Error().Msgf("Error occurred during terraform outputs syncing (refresh apply) %v", err)
		fmt.Println(stderr.String())
		return err
	}

	logger.Info().Msgf("Terraform outputs refreshed (synced) successfully for " + modulePath + " module")

	return nil

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
		return nil
	})

	if err != nil {
		logger.Error().Msgf("error walking the path %v\n", err)
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

func extractRefModuleFromString(input string) (string, error) {
	re, err := regexp.Compile(`getValueByKey\("([^"]+)"`)
	if err != nil {
		return "", err // Handle regex compilation error
	}
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1], nil // Return the captured group
	}
	return "", fmt.Errorf("no match found")
}

func getParentDirectory(filePath string) string {

	parent := filepath.Dir(filePath)

	return parent

}

func GetAllFilesInDir(dirPath string) []string {

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	files := []string{}
	for _, e := range entries {
		files = append(files, e.Name())
	}

	return files

}

func isValidTemplateFile(path string) bool {

	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	s := string(b)

	return strings.Contains(s, "getValueByKey")

}

func initTerraformDirectory(directory string) error {

	args := []string{"-chdir=" + directory, "init"}
	cmd := exec.Command("terraform", args...)

	// Set output to display in the console
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	err := cmd.Run()

	return err

}

func applyTerraform(directory string, autoApprove bool, init bool) error {

	TerraformTemplateProcessing(directory, true)

	if init {
		err := initTerraformDirectory(directory)
		if err != nil {
			return err
		}
	}

	// Construct the terraform apply command with -chdir
	args := []string{"-chdir=" + directory, "apply"}
	if autoApprove {
		args = append(args, "--auto-approve")
	}

	// Create the command with the constructed arguments
	cmd := exec.Command("terraform", args...)

	// Set output to display in the console
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error executing terraform apply: %w", err)
	}

	return nil
}

func copyDirectory(srcDir, dstDir string, excludePatterns []string) error {
	// Create the destination directory
	err := os.MkdirAll(dstDir, os.ModePerm)
	if err != nil {
		return err
	}

	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {

		// Modify files permissions
		args := []string{"-R", "777", dstDir}

		// Create the command with the constructed arguments
		cmd := exec.Command("chmod", args...)

		// Set output to display in the console
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Execute the command
		cmd.Run()

		if err != nil {
			return err
		}

		if shouldExclude(path, info, excludePatterns) {
			return nil
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, 0777)
		}
		return copyFile(path, dstPath)
	})

}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return dstFile.Sync()
}

func shouldExclude(path string, info os.FileInfo, excludePatterns []string) bool {
	if info.IsDir() {
		return false
	}
	for _, pattern := range excludePatterns {
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			return false
		}
		if matched {
			return true
		}
	}
	return false
}

func moveProjectToTemporaryDir(src string, dst string, excludedFiles []string) {

	err := copyDirectory(src, dst, excludedFiles)
	if err != nil {
		logger.Info().Msgf("Error copying directory: %v\n", err)
	} else {
		logger.Info().Msgf("Directory copied successfully.")
	}
}

func GetChildDirectory(filePath string) string {

	elements := []string{}

	elements = strings.Split(filePath, "/")

	childPath := elements[len(elements)-1]

	return childPath
}
