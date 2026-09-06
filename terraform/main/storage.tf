resource "aws_s3_bucket" "s3_storage" {
  bucket = "cloud-platform-kit-app-=${var.aws_account_id}"

}