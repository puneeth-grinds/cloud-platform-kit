output "ecr_api_gateway_url" {
  value       = aws_ecr_repository.ecr_api_gateway.repository_url
  description = "ECR Repository URL for api gateway"
}

output "ecr_vulnerability_scanner_url" {
  value       = aws_ecr_repository.ecr_vulnerability_scanner.repository_url
  description = "ECR Repository URL for Vulnerability Scanner"
}