package main

import (
	"fmt"
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

	structure := config.Structure

	srcDir := workdirPath
	tempDirPath := "/Users/vahagn/Documents/TerraformImprovmentProject/tmp"

	workdirPath = tempDirPath

	// Move project to temporary place to not interfere original files
	moveProjectToTemporaryDir(srcDir, tempDirPath, []string{"*.tfstate*", ".terraform"})
	logger.Debug().Msgf("Components map structure: %v", structure)

	dag, err := buildDAG(workdirPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	dag.Print(workdirPath)

	dot := dag.ToDot(workdirPath)
	err = ioutil.WriteFile("dag.dot", []byte(dot), 0644)
	if err != nil {
		fmt.Println("Error writing DOT file:", err)
		return
	}
	fmt.Println("DOT file written to dag.dot")

	for k, v := range dag.nodes {
		if len(v) == 0 {
			logger.Info().Msgf("Applying... %v", strings.TrimPrefix(k, workdirPath))
			// applyTerraform(k, true, true)
		} else {
			for _, dependentComponent := range v {
				logger.Info().Msgf("Applying... %v", strings.TrimPrefix(dependentComponent, workdirPath))
				// applyTerraform(dependentComponent, true, true)
			}

		}
	}

	//q := ApplyQueue{}

	//for _, rootComponent := range structure {
	//	for k := range rootComponent {
	//		logger.Debug().Msgf("%s -> %s ", k, rootComponent[k])
	//
	//		rootComponentPath := workdirPath + "/" + k
	//		logger.Debug().Msgf("Processing root components... %v", rootComponentPath)
	//
	//		for _, subComponentElement := range rootComponent[k] {
	//			logger.Debug().Msgf("Processing subcomponent... %v", subComponentElement)
	//
	//			subComponent := fmt.Sprintf("%v", subComponentElement)
	//
	//			// Getting files from each subdirectory and checking cross dependencies
	//			files := GetTerraformFiles(rootComponentPath + "/" + subComponent)
	//
	//			for _, v := range files {
	//				dependencies := GetDependency(v)
	//				if len(dependencies) != 0 {
	//					for _, v1 := range dependencies {
	//						refModule, _ := extractRefModuleFromString(v1)
	//
	//						q.Enqueue(refModule)
	//						q.Enqueue(subComponent)
	//					}
	//				}
	//			}
	//
	//		}
	//
	//	}
	//}

	//components := []string{}
	//for _, v := range q.elements {
	//	components = append(components, GetChildDirectory(v))
	//}
	//
	//logger.Info().Msgf("Pending apply: %v", strings.Join(components, ","))
	//
	//for _, v := range q.elements {
	//
	//	logger.Info().Msgf("In queue now: %v ", GetChildDirectory(v))
	//	//applyTerraform(v, true, true)
	//}

}
