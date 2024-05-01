module "ec2_instance" {
  source  = "terraform-aws-modules/ec2-instance/aws"

  name = "spot-instance"

  create_spot_instance = true
  instance_type          = "t2.micro"
  monitoring             = false
  subnet_id              = "subnet-0526c34347eadbc0b"
  #subnet_id              = "subnet-xyz"

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}
