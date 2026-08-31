# Backend state S3 bucket
output "state_bucket_name" {
  description = "Gives the bucket name on which tf state file is stored"
  value       = aws_s3_bucket.state_backend_tf.bucket
  sensitive   = false
}
