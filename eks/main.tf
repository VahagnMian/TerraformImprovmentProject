module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 19.0"

  cluster_name    = "dev-eks"
  cluster_version = "1.28"

  cluster_endpoint_public_access  = true

  vpc_id                   = getValueByKey("vpc", "vpc_id")
  subnet_ids               = getValueByKey("vpc", "public_subnet_ids")

  eks_managed_node_groups = {
    ondemand = {
      min_size     = 0
      max_size     = 1
      desired_size = 0

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