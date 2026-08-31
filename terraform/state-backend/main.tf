
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}
# Create S3 bucket
resource "aws_s3_bucket" "state_backend_tf" {
  bucket = "${var.project}-tf-state-${(data.aws_caller_identity.current.account_id)}"

  lifecycle {
    prevent_destroy = true
  }
}

# Enable versioning 
resource "aws_s3_bucket_versioning" "state_backend_tf_versioning" {
  bucket = aws_s3_bucket.state_backend_tf.id

  versioning_configuration {
    status = "Enabled"
  }

}

# Server-Side Encryption 
resource "aws_s3_bucket_server_side_encryption_configuration" "state_backend_tf_encryption" {
  bucket = aws_s3_bucket.state_backend_tf.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# Enable Public Access Protection
resource "aws_s3_bucket_public_access_block" "state_backend_tf_public_block" {
  bucket = aws_s3_bucket.state_backend_tf.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true

}
# Backend state to S3
terraform {
  backend "s3" {
    bucket       = "cloud-platform-kit-tf-state-512297269123"
    key          = "state-backend/terraform.tfstate"
    region       = "us-west-1"
    use_lockfile = true
    encrypt      = true
  }
}
