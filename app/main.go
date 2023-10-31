package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

type TerraformOutput struct {
	Sensitive bool          `json:"sensitive"`
	Type      []interface{} `json:"type"` // Use interface{} as type is a mixed structure.
	Value     []string      `json:"value"`
}

type TerraformOutputs map[string]TerraformOutput

func main() {
	// JSON data
	//jsonData := `{
	//    "private_subnet_ids": {
	//        "sensitive": false,
	//        "type": ["tuple", ["string", "string", "string"]],
	//        "value": ["subnet-0e630bed62b4f91a9", "subnet-0168fa5aba8893464", "subnet-0cdedca684b1a779a"]
	//    },
	//    "public_subnet_ids": {
	//        "sensitive": false,
	//        "type": ["tuple", ["string", "string", "string"]],
	//        "value": ["subnet-0e630bed62b4f91a9", "subnet-0168fa5aba8893464", "subnet-0cdedca684b1a779a"]
	//    }
	//}`

	outputData := string(getAllOutputs("../vpc"))

	var outputs TerraformOutputs
	err := json.Unmarshal([]byte(outputData), &outputs)
	if err != nil {
		log.Fatalf("Error parsing JSON: %s", err)
	}

	json.MarshalIndent(&outputs, "", "  ")

	fmt.Printf("Parsed Outputs: %+v\n", outputs)

}

func getAllOutputs(modulePath string) []byte {
	cmd := exec.Command("terraform", "output", "-json")
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
