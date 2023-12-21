module "ec2_instance" {
  source  = "terraform-aws-modules/ec2-instance/aws"

  name = "spot-instance"

  create_spot_instance = true
  instance_type          = "t2.micro"
  monitoring             = false
  subnet_id              = "subnet-0ed3c309930d86048"

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}
