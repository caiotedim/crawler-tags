output "vpc_id" {
  value = aws_vpc.vpc.id
}

output "igw_id" {
  value = aws_internet_gateway.internet_gateway.id
}

output "subnet" {
  value = aws_subnet.subnet.*.id
}