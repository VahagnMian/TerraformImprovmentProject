package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

func TerraformTemplateProcessing(directory string, inputFileName string, overwriteTF bool) {

	logger := initLogger()

	destinationTemplateFile := directory + "/" + inputFileName

	file, err := os.Open(destinationTemplateFile)
	defer file.Close()

	if err != nil {
		log.Fatal("[ Error ] error occured during reading of: ", destinationTemplateFile, " please check if file exists, and has right permissions")
	}
	logger.Println("File " + destinationTemplateFile + " read successfully")

	pattern := regexp.MustCompile(`getValueByKey\("([^"]+)", "([^"]+)"\)`)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		matches := pattern.FindStringSubmatch(line)

		if matches != nil {

			outputs := getAllOutputs("../"+matches[1], false)

			actualValue := getValueByKey(matches[2], parseHCL(string(outputs)))

			line = pattern.ReplaceAllString(line, actualValue)
			logger.Println("Successfully replaced functions with their actual values")

		}

		fmt.Println(line)
		writeResultToFile(line, appendProcessedToTf("../ec2/main.tf"), overwriteTF)

	}

	renameFile(overwriteTF, appendProcessedToTf("../ec2/main.tf"))

}

func writeResultToFile(line string, filePath string, overwriteTF bool) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening/creating file: ", err)
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, line)
	if err != nil {
		log.Fatal("Error writing to file: ", err)
	}
}

func renameFile(overwriteTF bool, filePath string) {

	if overwriteTF {
		err1 := os.Rename(filePath, trimProcessedFromTf(filePath))
		if err1 != nil {
			fmt.Println(err1)
		}
	}
}
