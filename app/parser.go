package main

import (
	"bufio"
	"fmt"
	logger "github.com/rs/zerolog/log"
	"log"
	"os"
	"regexp"
)

// Method is responsible for overwriting the terraform "template" file with the actuall passed values
func TerraformTemplateProcessing(directory string, inputFileName string, overwriteTF bool) {

	destinationTemplateFile := directory + "/" + inputFileName

	// Read file
	file, err := os.Open("../" + destinationTemplateFile)
	defer file.Close()

	// Error handling
	if err != nil {
		log.Fatal("[ Error ] error occured during reading of: ", destinationTemplateFile, " please check if file exists, and has right permissions")
		return
	}
	logger.Info().Msgf("File " + destinationTemplateFile + " read successfully")

	pattern := regexp.MustCompile(`getValueByKey\("([^"]+)", "([^"]+)"\)`)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		matches := pattern.FindStringSubmatch(line)

		if matches != nil {

			outputs := getAllOutputs("../"+matches[1], false)

			actualValue := getValueByKey(matches[2], parseHCL(string(outputs)))

			line = pattern.ReplaceAllString(line, actualValue)
			logger.Info().Msgf("Successfully replaced functions with their actual values")

		}

		writeResultToFile(line, appendProcessedToTf("../"+directory+"/main.tf"))

	}

	renameFile(overwriteTF, appendProcessedToTf("../"+directory+"/main.tf"))

}

func writeResultToFile(line string, filePath string) {
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
