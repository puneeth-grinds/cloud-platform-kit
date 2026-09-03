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
  security_group_ids  = [aws_security_group.allow_tls.id]
}

resource "aws_vpc_endpoint" "ecr_interface_endpoint_api" {
  vpc_id              = aws_vpc.vpc.id
  service_name        = "com.amazonaws.${var.aws_region}.ecr.api"
  vpc_endpoint_type   = "Interface"
  subnet_ids          = [aws_subnet.private_1a.id, aws_subnet.private_1b.id]
  private_dns_enabled = true
  security_group_ids  = [aws_security_group.allow_tls.id]
}

resource "aws_security_group" "allow_tls" {
  name        = "allow_tls"
  description = "Allow TLS inbound traffic"
  vpc_id      = aws_vpc.vpc.id

  tags = {
    Name = "security_group_allow_https"
  }
}

resource "aws_vpc_security_group_ingress_rule" "allow_tls_igress_rule" {
  security_group_id = aws_security_group.allow_tls.id

  # Restricted source network
  cidr_ipv4   = "10.0.0.0/16"
  from_port   = 443
  ip_protocol = "tcp"
  to_port     = 443

}