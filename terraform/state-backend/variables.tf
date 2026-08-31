# AWS account ID
variable "aws_account_id" {
  type        = string
  description = "Personal AWS account ID"
}

# AWS Default Region
variable "aws_region" {
  type        = string
  description = "Specifies the AWS region"
  default     = "us-west-1"
}
# Project name
variable "project" {
  type        = string
  description = "Specifies the project name"
  default     = "cloud-platform-kit"
}