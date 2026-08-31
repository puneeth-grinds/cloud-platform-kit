resource "aws_s3_bucket" "state_backend_tf" {
  bucket = "state-backend-tf"
}

resource "aws_s3_bucket_versioning" "state_backend_tf_versioning" {
  bucket = aws_s3_bucket.state_backend_tf.id

  versioning_configuration {
    status = "Enabled"
  }

}