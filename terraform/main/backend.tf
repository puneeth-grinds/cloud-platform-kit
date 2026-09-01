terraform {
  backend "s3" {
    bucket       = "cloud-platform-kit-tf-state-512297269123"
    key          = "main/terraform.tfstate"
    region       = "us-west-1"
    encrypt      = true
    use_lockfile = true
  }
}