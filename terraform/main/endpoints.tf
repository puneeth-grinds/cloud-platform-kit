resource "aws_vpc_endpoint" "s3_gateway_endpoint" {
  vpc_id            = aws_vpc.vpc.id
  service_name      = "com.amazonaws.${var.aws_region}.s3"
  vpc_endpoint_type = "Gateway"
  route_table_ids   = [aws_route_table.private-route-table.id]
}

resource "aws_vpc_endpoint" "ecr_interface_endpoint_dkr" {
  vpc_id              = aws_vpc.vpc.id
  service_name        = "com.amazonaws.${var.aws_region}.ecr.dkr"
  vpc_endpoint_type   = "Interface"
  subnet_ids          = [aws_subnet.private_1a.id, aws_subnet.private_1b.id]
  private_dns_enabled = true
}

resource "aws_vpc_endpoint" "ecr_interface_endpoint_api" {
  vpc_id              = aws_vpc.vpc.id
  service_name        = "com.amazonaws.${var.aws_region}.ecr.api"
  vpc_endpoint_type   = "Interface"
  subnet_ids          = [aws_subnet.private_1a.id, aws_subnet.private_1b.id]
  private_dns_enabled = true
}