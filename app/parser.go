package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

func TerraformTemplateProcessing(directory string, inputFileName string) {

	destinationTemplateFile := directory + "/" + inputFileName

	file, err := os.Open(destinationTemplateFile)

	if err != nil {
		log.Fatal("[ Error ] error occured during reading of: ", destinationTemplateFile, " please check if file exists, and has right permissions")
	}

	defer file.Close()

	// Defining regexp pattern
	pattern := regexp.MustCompile(`getValueByKey\("([^"]+)", "([^"]+)"\)`)

	// Scanning file line by line
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		matches := pattern.FindStringSubmatch(line)

		if matches != nil {

			outputs := getAllOutputs("../"+matches[1], false)

			actualValue := getValueByKey(matches[2], parseHCL(string(outputs)))

			//fmt.Println(actualValue)
			line = pattern.ReplaceAllString(line, actualValue)
		}

		fmt.Println(line)
	}

}
