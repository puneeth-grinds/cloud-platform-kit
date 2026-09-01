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
    condition     = length(var.aws_account_id) == 12
    error_message = "AWS account id must be valid of 12 characters"
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