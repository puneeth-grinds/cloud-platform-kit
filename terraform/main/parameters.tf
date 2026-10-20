resource "aws_ssm_parameter" "api_gateway_api_key_ssm" {
  name        = "/cloud-platform-kit/api-gateway/api-key"
  description = "Stores the API key for API Gateway"
  type        = "SecureString"
  value       = var.api_gateway_api_key

  tags = {
    Component = "api-gateway"
  }
}

resource "aws_ssm_parameter" "grafana_otlp_endpoint" {
  name        = "/cloud-platform-kit/grafana/otlp-endpoint"
  description = "Stores the Grafana Cloud OTLP endpoint"
  type        = "SecureString"
  value       = var.grafana_otlp_endpoint

  tags = {
    Component = "observability"
  }
}

resource "aws_ssm_parameter" "grafana_otlp_headers" {
  name        = "/cloud-platform-kit/grafana/otlp-headers"
  description = "Stores the Grafana Cloud OTLP authorization header"
  type        = "SecureString"
  value       = var.grafana_otlp_headers

  tags = {
    Component = "observability"
  }
}

