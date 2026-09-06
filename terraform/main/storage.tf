resource "aws_s3_bucket" "s3_storage" {
  bucket = "cloud-platform-kit-app-=${var.aws_account_id}"
  tags = {
    Name = "cloud-platform-kit-storage-bucket"
  }
}

resource "aws_s3_bucket_versioning" "s3_storage_versioning" {
  bucket = aws_s3_bucket.s3_storage.id
  versioning_configuration {
    status = "Enabled"
  }
}