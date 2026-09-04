resource "aws_ecr_repository" "ecr" {
  name                 = "ecr for api gateway"
  image_tag_mutability = "IMMUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}