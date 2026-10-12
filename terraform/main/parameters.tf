resource "aws_ssm_parameter" "api_gateway_api_key_ssm" {
  name = "api_gatewayapi_key"
  description = "Stores the API key for API Gateway"
  type = "SecureString"
  value = var.api_gateway_api_key
  

}