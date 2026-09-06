
resource "aws_db_subnet_group" "rds_subnet_group" {
  name       = "RDS_Subnet"
  subnet_ids = [aws_subnet.private_1a, aws_subnet.private_1b]

  tags = {
    Name = "RDS subnet group"
  }
}