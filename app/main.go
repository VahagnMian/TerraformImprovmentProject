package main

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"log"
	"os/exec"
)

func main() {

	outputs := getAllOutputs("../vpc", false)

	fmt.Println(string(outputs))

	getValueByKey("private_subnet_ids", parseHCL(string(outputs)))

}

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

func getValueFrom(input string) string {
	return "Processed Value: " + input
}

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

func getValueByKey(key string, result map[string]interface{}) {
	fmt.Println(result[key])
}
