provider "aws" {
  region              = "us-west-1"
  allowed_account_ids = [var.aws_account_id]
}