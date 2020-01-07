
locals {
  template_count = lookup(var.cloud_init, "template", null) == null ? 0 : var.instance_count
  template_vars  = lookup(var.cloud_init, "vars", {})
  eip_count      = var.associate_eip ? var.instance_count : 0
  instance_ports = [
    for port in var.ports :
    [for i in range(var.instance_count) :
      port
    ]
  ]
}

data "template_file" "cloud_init" {
  count = local.template_count
  template = templatefile(var.cloud_init.template,
    merge(
      local.template_vars,
      { index = count.index },
      { volume = var.volume },
      { tags = merge(var.tags, { Name = format("%s-%d", var.name, count.index + 1) }) }
    )
  )
}


resource "aws_network_interface" "port0" {
  count             = length(local.instance_ports) > 0 ? length(local.instance_ports[0]) : 0
  description       = lookup(local.instance_ports[0][count.index], "description", null)
  private_ips       = lookup(local.instance_ports[0][count.index], "private_ips", null)
  private_ips_count = lookup(local.instance_ports[0][count.index], "private_ips_count", null)
  subnet_id         = lookup(local.instance_ports[0][count.index], "subnet_id", null)
  source_dest_check = lookup(local.instance_ports[0][count.index], "source_dest_check", null)
  security_groups   = lookup(local.instance_ports[0][count.index], "security_groups", null)
  tags = merge(
    var.tags,
    { Name = format("%s-%d", var.name, count.index + 1) }
  )
}

resource "aws_instance" "ec2" {
  count                  = var.instance_count
  ami                    = var.ami
  subnet_id              = var.subnet_id
  instance_type          = var.instance_type
  user_data              = element(data.template_file.cloud_init.*.rendered, count.index)
  ipv6_address_count     = length(local.instance_ports) == 0 ? var.ipv6_address_count : null
  ipv6_addresses         = length(local.instance_ports) == 0 ? var.ipv6_addresses : null
  source_dest_check      = length(local.instance_ports) == 0 ? var.source_dest_check : null
  iam_instance_profile   = var.iam_instance_profile
  vpc_security_group_ids = var.vpc_security_group_ids
  associate_public_ip_address = var.associate_public_ip_address

  dynamic "network_interface" {
    iterator = port
    for_each = length(aws_network_interface.port0.*.id) > 0 ? [aws_network_interface.port0[count.index].id] : []
    content {
      device_index         = port.key
      network_interface_id = port.value
    }

  }

  ebs_optimized = var.ebs_optimized

  volume_tags = merge(
    var.tags,
    var.volume_tags
  )

  root_block_device {
    volume_size           = lookup(var.root_block_device, "volume_size", null)
    volume_type           = lookup(var.root_block_device, "volume_type", null)
    delete_on_termination = lookup(var.root_block_device, "delete_on_termination", null)
  }

  disable_api_termination              = var.disable_api_termination
  instance_initiated_shutdown_behavior = var.instance_initiated_shutdown_behavior
  placement_group                      = var.placement_group
  tenancy                              = var.tenancy

  tags = merge(
    var.tags,
    { Name = format("%s-%d", var.name, count.index + 1) }
  )
}

resource "aws_route_table" "default" {
  vpc_id = var.vpc_id

  tags = merge({
    Name    = format("default_table_%s", var.zone),
    purpose = "common_services"
    },
    var.tags
  )

  lifecycle {
    ignore_changes = all
  }

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = var.igw_id
  }
  route {
    cidr_block           = "10.0.0.0/8"
    network_interface_id = element(aws_network_interface.port0.*.id, 0)
  }
}

resource "aws_route_table_association" "default" {
  subnet_id      = lookup(local.instance_ports[0][0], "subnet_id", null)
  route_table_id = aws_route_table.default.id
}

resource "aws_eip" "eip" {
  count    = local.eip_count
  instance = element(aws_instance.ec2.*.id, count.index)
  vpc      = true
  tags = merge(
    var.tags,
    { Name = format("%s-%d", var.name, count.index + 1) }
  )
}

resource "aws_ebs_volume" "volume" {
  count             = lookup(var.volume, "size", null) == null ? 0 : var.instance_count
  size              = var.volume.size
  availability_zone = var.tags.zone
  type              = lookup(var.volume, "type", "gp2")

  tags = merge(
    var.tags,
    { Name = format("%s-%d", lookup(var.volume, "Name", var.name), count.index + 1) }
  )
}

resource "aws_volume_attachment" "ebs_attachment" {
  count       = lookup(var.volume, "size", null) == null ? 0 : var.instance_count
  device_name = var.volume.device
  volume_id   = element(aws_ebs_volume.volume.*.id, count.index)
  instance_id = element(aws_instance.ec2.*.id, count.index)
  skip_destroy = var.skip_destroy
}
