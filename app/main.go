package main

func main() {

	refreshTerraformOutputs("vpc")

	TerraformTemplateProcessing("../eks", "main.tf")

}
