module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 19.0"

  cluster_name    = "dev-eks"
  cluster_version = "1.28"

  cluster_endpoint_public_access  = true

  #vpc_id                   = "vpc-0eb0cf447f3788272"

  vpc_id                   = getValueFrom(vpc.vpc_id)
  #subnet_ids               = ["subnet-0e630bed62b4f91a9", "subnet-0168fa5aba8893464", "subnet-0cdedca684b1a779a"]
  subnet_ids               = getValueFrom(vpc.private_subnets)

  eks_managed_node_groups = {
    ondemand = {
      min_size     = 1
      max_size     = 1
      desired_size = 1

      instance_types = ["t2.micro"]
      capacity_type  = "SPOT"
    }
  }

  # aws-auth configmap
  manage_aws_auth_configmap = true

  aws_auth_roles = []

  aws_auth_users = []

  aws_auth_accounts = []

  tags = {
    Environment = "dev"
    Terraform   = "true"
  }
}