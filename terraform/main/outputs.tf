output "ecr_api_gateway_url" {
  value       = aws_ecr_repository.ecr_api_gateway.repository_url
  description = "ECR Repository URL for api gateway"
}

output "ecr_vulnerability_scanner_url" {
  value       = aws_ecr_repository.ecr_vulnerability_scanner.repository_url
  description = "ECR Repository URL for Vulnerability Scanner"
}

output "s3_bucket_storage_name" {
  value       = aws_s3_bucket.s3_storage.bucket
  description = "S3 storage bucket name"
}
output "s3_bucket_storage_arn" {
  value       = aws_s3_bucket.s3_storage.arn
  description = "S3 storage bucket arn"
}

output "ecs_task_execution_role_arn" {
  value       = aws_iam_role.ecs_task_execution.arn
  description = "ECS task execution role ARN"
}

output "api_gateway_task_role_arn" {
  value       = aws_iam_role.api_gateway_task.arn
  description = "API Gateway ECS task role ARN"
}

output "vulnerability_scanner_task_role_arn" {
  value       = aws_iam_role.vulnerability_scanner_task.arn
  description = "Vulnerability scanner ECS task role ARN"
}

output "alb_dns_name" {
  value       = aws_lb.lb.dns_name
  description = "Provides the lb dns name"
}