
locals {
  all_vars = yamldecode(file("vars.yaml"))
  all_tags = lookup(local.all_vars, "tags", {} )
  vpc_blocks_count = lookup(local.all_vars.vpc, "blocks", 0 ) != 0 ? length(local.all_vars.vpc.blocks) : 0
  vpc_subnet_count = lookup(local.all_vars.vpc, "subnet", 0) != 0 ? length(local.all_vars.vpc.subnet) : 0
  vpc_blocks =  lookup(local.all_vars.vpc, "blocks", [] )
  tags = merge(
    local.all_tags,
    {zone = local.all_vars.zone},
  )
  instances_template = "templates/cloud-config.tpl"
}

provider "aws" {
  region = local.all_vars.region
}

terraform {
  backend "local" {
    path = "./terraform.tfstate"
  }
}

resource "aws_vpc" "vpc" {
  cidr_block = local.all_vars.vpc.main_block
  tags = local.tags
}

resource "aws_internet_gateway" "internet_gateway" {
  vpc_id = aws_vpc.vpc.id

  tags = local.tags
}

resource "aws_subnet" "subnet" {
  count                   = local.vpc_subnet_count
  vpc_id                  = aws_vpc.vpc.id
  cidr_block              = local.all_vars.vpc.subnet[count.index]
  map_public_ip_on_launch = false
  availability_zone       = local.all_vars.zone

  tags = merge({
    Name = format("crawler_tags_%s", local.all_vars.zone) },
    local.tags
  )
}

module "crawler_tags_instances" {
  source          = "./modules/aws/instances"
  name            = "crawler-tags"
  instance_count  = local.all_vars.instance.count
  vpc_id          = aws_vpc.vpc.id
  igw_id          = aws_internet_gateway.internet_gateway.id
  associate_eip   = true
  ports = [
    {
      subnet_id         = aws_subnet.subnet[0].id
      security_groups   = [module.security_group.id]
      source_dest_check = false
      tags              = local.tags
    }
  ]
  volume          = lookup(local.all_vars.instance, "volume", {})
  ami             = local.all_vars.instance.ami
  instance_type   = local.all_vars.instance.type
  tags            = local.tags
  cloud_init = {
    template = local.instances_template,
    vars     = { ssh_keys : local.all_vars.ssh_keys }
  }
  zone       = local.all_vars.zone
}
