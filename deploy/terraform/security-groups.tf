module "security_group" {
  source  = "./modules/aws/security_groups"
  name    = "craler_tags_sg"
  tags    = local.tags
  vpc_id  = aws_vpc.vpc.id
  cidr    = local.all_vars.vpc.subnet
}
