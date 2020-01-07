variable "vpc_id" {
  description = "vpc_id to where the security group will be created"
  type        = string
}

variable "description" {
  type        = string
  description = "security group description"
  default     = null
}

variable "name" {
  description = "security group name"
  type        = string
}

variable "tags" {
  description = "security group tags"
  default     = {}
}

variable "cidr" {
  default     = []
}