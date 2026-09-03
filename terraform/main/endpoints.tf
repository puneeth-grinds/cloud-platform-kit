resource "aws_vpc_endpoint" "s3_gateway_endpoint" {
  region = var.aws_region
  vpc_id = aws_vpc.vpc.id
  service_name = "com.amazonaws.${var.aws_region}.s3"
  route_table_ids = aws_route_table.private-route-table
}