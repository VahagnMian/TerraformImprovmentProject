package main

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"log"
	"os"
	"os/exec"
	"strings"
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

	logger := initLogger()
	logger.Println("Starting Terraform outputs syncing")

	_, err := exec.Command("terraform", "-chdir=../"+modulePath, "apply", "-refresh-only", "-auto-approve").Output()

	if err != nil {
		log.Fatal("Error accured during terraform outputs syncing (refrsh apply) ", err)
	}
	logger.Println("Terraform outputs refreshed (synced) successfully for " + modulePath + " module")

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

func initLogger() *log.Logger {
	return log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
}

//func writeResultToFile(line string, filePath string) *bufio.Writer {
//	logger := initLogger()
//
//	_, err := os.Stat(filePath)
//
//	var file *os.File
//	var a = 1
//
//	if err != nil {
//		logger.Println("File ", filePath, "doesn't exists, creating it now!")
//		file, _ = os.Create(filePath)
//		//if err1 != nil {
//		//	logger.Fatal("Error creating ", filePath, ": ", err1)
//		//}
//
//	}
//
//	w := bufio.NewWriter(f)
//
//	w.WriteString(line)
//
//	return w
//
//}

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
