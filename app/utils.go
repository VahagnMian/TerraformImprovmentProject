package main

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"log"
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

	//if isHCLString(value) {
	//
	//}

	if isHCLArray(value) {
		value = makeHCLArrayFromHCLOutput(value)
	}

	return value

}

// refreshTerraformOutputs is used to sync terraform code and state so code can retrieve outputs
func refreshTerraformOutputs(modulePath string) {

	out, err := exec.Command("terraform", "-chdir=../"+modulePath, "apply", "-refresh-only", "-auto-approve").Output()

	fmt.Println(string(out))

	//fmt.Println(out)

	if err != nil {
		log.Fatal(err)
	}

}

// isHCLArray checking if terraform output is array, so I can convert it to terraform friendly array
func isHCLArray(s string) bool {
	return strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]")
}

// makeHCLArrayFromHCLOutput this function used to convert HCL output formated Array to terraform friendly array
func makeHCLArrayFromHCLOutput(s string) string {

	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")

	elements := strings.Fields(s)

	var hclArray string

	for i := 0; i < len(elements); i++ {
		hclArray = hclArray + elements[i] + ", "
	}
	hclArray = strings.TrimSuffix(hclArray, ", ")
	hclArray = "[" + hclArray + "]"

	return hclArray

}
