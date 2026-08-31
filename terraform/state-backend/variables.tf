variable "aws_account_id" {
  type        = string
  description = "Personal AWS account ID"
}

variable "aws_region" {
  type        = string
  description = "Specifies the AWS region"
  default     = "eu-east-1"
}

variable "project" {
  type        = string
  description = "Specifies the project name"
  default     = "cloud-platform-kit"
}