resource "aws_ecr_repository" "ecr_api_gateway" {
  name                 = "ecr for api gateway"
  image_tag_mutability = "IMMUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_repository" "ecr_vulnerability_scanner" {
  name                 = "ecr for vulnerability scanner"
  image_tag_mutability = "IMMUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}