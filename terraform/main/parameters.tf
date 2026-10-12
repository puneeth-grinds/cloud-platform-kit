resource "aws_ssm_parameter" "api_gateway_api_key_ssm" {
  name        = "/cloud-platform-kit/api-gateway/api-key"
  description = "Stores the API key for API Gateway"
  type        = "SecureString"
  value       = var.api_gateway_api_key

  tags = {
    Component = "api-gateway"
  }
}
