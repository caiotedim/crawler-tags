variable "ports" {
  default = []
}

variable "name" {
  default = "instance_name"
}
variable "vpc_id" {
  description = "vpc id that is going to be used in this deployment"
  default     = ""
}

variable "instance_count" {
  description = "Count of instances to be created by terraform"
  default     = "1"
}


variable "ami" {
  description = "ID of AMI to use for the instance"
}

variable "placement_group" {
  description = "The Placement Group to start the instance in"
  default     = ""
}

variable "tenancy" {
  description = "The tenancy of the instance (if the instance is running in a VPC). Available values: default, dedicated, host."
  default     = "default"
}

variable "ebs_optimized" {
  description = "If true, the launched EC2 instance will be EBS-optimized"
  default     = false
}

variable "disable_api_termination" {
  description = "If true, enables EC2 Instance Termination Protection"
  default     = false
}

variable "instance_initiated_shutdown_behavior" {
  description = "Shutdown behavior for the instance"
  default     = ""
}

variable "instance_type" {
  description = "Instance Type to be created based in AWS flavors size"
  default     = "t2.nano"
}

variable "key_name" {
  description = "The key name to use for the instance"
  default     = ""
}

variable "monitoring" {
  description = "If true, the launched EC2 instance will have detailed monitoring enabled"
  default     = false
}

variable "vpc_security_group_ids" {
  description = "A list of security group IDs to associate with"
  type        = list
  default     = null
}

variable "security_groups" {
  description = "A list of security group IDs to associate with"
  default     = ["default"]
}

variable "subnet_id" {
  description = "The VPC Subnet ID to launch in"
  default     = null
}

variable "associate_public_ip_address" {
  description = "If true, the EC2 instance will have associated public IP address"
  default     = null
}

variable "associate_eip" {
  description = "If true, the EC2 instance will have associated EIP address"
  default     = false
}

variable "private_ip" {
  description = "Private IP address to associate with the instance in a VPC"
  default     = ""
}

variable "source_dest_check" {
  description = "Controls if traffic is routed to the instance when the destination address does not match the instance. Used for NAT or VPNs."
  default     = true
}

variable "cloud_init" {
  description = "The user data to provide when launching the instance"
  default     = ""
}

variable "iam_instance_profile" {
  description = "The IAM Instance Profile to launch the instance with. Specified as the name of the Instance Profile."
  default     = ""
}

variable "ipv6_address_count" {
  description = "A number of IPv6 addresses to associate with the primary network interface. Amazon EC2 chooses the IPv6 addresses from the range of your subnet."
  default     = 0
}

variable "ipv6_addresses" {
  description = "Specify one or more IPv6 addresses from the range of the subnet to associate with the primary network interface"
  default     = []
}

variable "tags" {
  description = "A mapping of tags to assign to the resource"
  default     = {}
}

variable "volume_tags" {
  description = "A mapping of tags to assign to the devices created by the instance at launch time"
  default     = {}
}

variable "root_block_device" {
  description = "Customize details about the root block device of the instance. See Block Devices below for details"
  default     = {}
}

variable "ebs_block_device" {
  description = "Additional EBS block devices to attach to the instance"
  default     = []
}

variable "ephemeral_block_device" {
  description = "Customize Ephemeral (also known as Instance Store) volumes on the instance"
  default     = []
}

variable "volume" {
  description = "Create volume with var.volume_size and attach to instance"
  default     = {}
}

variable "skip_destroy" {
  default = true
}
variable "device_name" {
  description = "Device name for volume"
  default     = ""
}

variable "igw_id" {
  default = ""
}

variable "zone" {
  default = ""
}
