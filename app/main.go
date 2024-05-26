package main

import (
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Config struct {
	WorkdirPath string                     `yaml:"workdir"`
	Structure   []map[string][]interface{} `yaml:"structure"`
}

var config Config

func main() {
	logger.Logger = logger.Output(zerolog.ConsoleWriter{Out: colorable.NewColorableStdout()})

	// Reading the config file and unmarshalling
	data, err := ioutil.ReadFile("../vinfra.yaml")
	checkErr(err, "Error reading file")

	err = yaml.Unmarshal(data, &config)
	checkErr(err, "Error unmarshaling YAML")

	logger.Debug().Msgf(config.WorkdirPath)
	createInfraStructure(config.WorkdirPath)

	workdirPath := config.WorkdirPath

	srcDir := workdirPath
	tempDirPath := "/Users/vahagn/Documents/TerraformImprovmentProject/tmp"

	workdirPath = tempDirPath

	moveProjectToTemporaryDir(srcDir, tempDirPath, []string{"*.tfstate*", ".terraform"})

	q := ApplyQueue{}

	files := GetTerraformFiles(workdirPath)

	logger.Info().Msgf("All terraform configuration files found")
	for _, file := range files {

		logger.Info().Msgf("%v", file)
	}

	for _, v := range files {
		getParentDirectory(v)

	}

	// Going through all the files in specified workdir and developing queue
	componentsApply := make(map[string]bool)
	for _, v := range files {
		dependencies := GetDependency(v)
		if len(dependencies) != 0 {
			for k, v1 := range dependencies {
				refModule, _ := extractRefModuleFromString(v1)
				mainModule := getParentDirectory(k)
				refModule = getParentDirectory(mainModule) + "/" + refModule

				q.Enqueue(refModule)
				q.Enqueue(mainModule)
			}
		} else {
			componentsApply[getParentDirectory(v)] = true
			q.Enqueue(getParentDirectory(v))

		}
	}

	components := []string{}
	for _, v := range q.elements {
		components = append(components, GetChildDirectory(v))
	}

	logger.Info().Msgf("Pending apply: %v", strings.Join(components, ","))

	for _, v := range q.elements {

		logger.Info().Msgf("In queue now: %v ", GetChildDirectory(v))
		//applyTerraform(v, true, true)
	}

}
