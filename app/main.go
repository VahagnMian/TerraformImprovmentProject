package main

func main() {

	refreshTerraformOutputs("vpc")
	//refreshTerraformOutputs("vpc")
	TerraformTemplateProcessing("../ec2", "main.tf", false)

}
