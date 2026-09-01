variable "aws_region" {
  type        = string
  description = "Specifies the AWS region"
  default     = "us-west-1"
}

variable "aws_account_id" {
  type        = string
  description = "Personal AWS account ID"
  default     = "512297269123"

  validation {
    condition     = can(regex("^[0-9]{12}$", var.aws_account_id))
    error_message = "AWS account id must be valid of 12 digits"
  }
}

variable "project" {
  type        = string
  description = "Specifies the project name"
  default     = "cloud-platform-kit"
}

variable "environment" {
  type        = string
  description = "Specifies the default environment"
  default     = "dev"
}

variable "vpc_cidr_block" {
  type        = string
  description = "Specifies the VPC CIDR Block"
}
variable "public_1a" {
  type        = string
  description = "Subnet for public_1a"
}
variable "public_1b" {
  type        = string
  description = "Subnet for public_1b"
}
variable "privatec_1a" {
  type        = string
  description = "Subnet for privatec_1a"
}
variable "private_1b" {
  type        = string
  description = "Subnet for private_1b"
}