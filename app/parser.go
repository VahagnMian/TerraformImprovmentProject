package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

func TerraformTemplateProcessing(directory string, inputFileName string) {

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

		writeResultToFile(line, appendProcessedToTf("../ec2/main.tf"))

	}

}

func writeResultToFile(line string, filePath string) {
	logger := initLogger()

	// Determine file mode (create or append)
	mode := os.O_WRONLY
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		mode |= os.O_CREATE
	} else {
		mode |= os.O_APPEND
	}

	// Open or create the file
	file, err := os.OpenFile(filePath, mode, 0644)
	if err != nil {
		logger.Fatal("Error opening/creating file: ", err)
	}
	defer file.Close()

	// Write the line to the file
	_, err = fmt.Fprintln(file, line)
	if err != nil {
		logger.Fatal("Error writing to file: ", err)
	}
}
