# Create S3 bucket
resource "aws_s3_bucket" "state_backend_tf" {
  bucket = "state-backend-tf"
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