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

resource "aws_ecr_lifecycle_policy" "ecr_policy_api_gateway_p1" {
  repository = aws_ecr_repository.ecr_api_gateway.name

  policy = <<EOF
{
  "rules": [
    {
      "rulePriority": 1,
      "description": "Expire images older than 1 days",
      "selection": {
        "tagStatus": "untagged",
        "countType": "sinceImagePushed",
        "countUnit": "days",
        "countNumber": 1
      },
      "action": {
        "type": "expire"
      }
    }
  ]
}
EOF

}

resource "aws_ecr_lifecycle_policy" "ecr_policy_api_gateway_p2" {
  repository = aws_ecr_repository.ecr_api_gateway.name

  policy = <<EOF
{
  "rules": [
    {
      "rulePriority": 1,
      "description": "Keep only the newest image ",
      "selection": {
        "tagStatus": "untagged",
        "countType": "sinceImagePushed",
        "countUnit": "days",
        "countNumber": 10
      },
      "action": {
        "type": "expire"
      }
    }
  ]
}
EOF

}

resource "aws_ecr_lifecycle_policy" "ecr_policy_vul_scanner_p1" {
  repository = aws_ecr_repository.ecr_vulnerability_scanner.name

  policy = <<EOF
{
  "rules": [
    {
      "rulePriority": 1,
      "description": "Expire images older than 1 days",
      "selection": {
        "tagStatus": "untagged",
        "countType": "sinceImagePushed",
        "countUnit": "days",
        "countNumber": 1
      },
      "action": {
        "type": "expire"
      }
    }
  ]
}
EOF

}