module "ec2_instance" {
  source  = "terraform-aws-modules/ec2-instance/aws"

  name = "spot-instance"

  create_spot_instance = true
  instance_type          = "t2.micro"
  monitoring             = false
  subnet_id              = getValueByKey("vpc", "public_first_subnet_id")
  #subnet_id              = "subnet-xyz"

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}