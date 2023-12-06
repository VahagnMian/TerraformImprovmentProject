#output "vpc_id" {
#  value = module.vpc.vpc_id
#}

output "private_subnet_ids" {
  value = module.vpc.private_subnets
}

output "public_subnet_ids" {
  value = module.vpc.private_subnets
}

output "public_first_subnet_id" {
  value = module.vpc.private_subnets[0]
}

output "vpc_id" {
  value = module.vpc.vpc_id
}