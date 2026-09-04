output "ecr_api_gateway_url" {
  value = aws_ecr_repository.ecr_api_gateway.repository_url
  description = "ECR Repository URL for api gateway"
}

