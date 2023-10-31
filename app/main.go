package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	//// Open the template file
	//file, err := os.Open("../eks/main.tf")
	//if err != nil {
	//	panic(err)
	//}
	//defer file.Close()
	//
	//// Regular expression to match getValueFrom(...) pattern
	//regex := regexp.MustCompile(`getValueFrom\(([^)]+)\)`)
	//
	//scanner := bufio.NewScanner(file)
	//for scanner.Scan() {
	//	line := scanner.Text()
	//
	//	// Find all matches in the line
	//	matches := regex.FindAllStringSubmatch(line, -1)
	//
	//	for _, match := range matches {
	//		// match[0] is the whole match, match[1] is the captured group
	//		fmt.Println("Match: ", match)
	//		fmt.Println("Match[0]: ", match[0])
	//		fmt.Println("Match[1]: ", match[1])
	//		//replacement := getValueFrom(match[1])
	//
	//		// Replace the matched string with its processed value
	//		//line = strings.Replace(line, match[0], replacement, 1)
	//	}
	//
	//	fmt.Println(line)
	//}
	//
	//if err := scanner.Err(); err != nil {
	//	panic(err)
	//}

	outputVPC := getAllOutputs("../vpc/")

	e := json.Unmarshal(outputVPC, &c)

	if e != nil {
		log.Panicf("Error happend:", e)
	}
}

func getValueFrom(input string) string {
	return "Processed Value: " + input
}

func getAllOutputs(modulePath string) []byte {
	cmd := exec.Command("terraform", "output", "-json")
	cmd.Dir = modulePath

	// Run the command
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing Terraform:", err)
	}

	return output
}
