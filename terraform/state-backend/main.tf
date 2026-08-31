terraform {
  backend "s3" {
    bucket       = "state-backend-tf"
    key          = "environments/dev/terraform.tfstate"
    region       = var.aws_region
    encrypt      = true
    use_lockfile = true
  }
}