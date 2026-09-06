
resource "aws_db_subnet_group" "rds_subnet_group" {
  name       = "rds"
  subnet_ids = [aws_subnet.private_1a.id, aws_subnet.private_1b.id]

  tags = {
    Name = "RDS subnet group"
  }
}

resource "aws_db_instance" "RDS_DB" {
  engine                 = "PostgreSQL"
  engine_version         = "16"
  instance_class         = "db.t3.micro"
  allocated_storage      = 20
  multi_az               = false
  availability_zone      = var.availability_zone_1a
  db_name                = var.rds_database_name
  username               = var.rds_database_username
  password               = var.rds_database_password
  publicly_accessible    = false
  skip_final_snapshot    = true
  deletion_protection    = false
  vpc_security_group_ids = [aws_security_group.rds_sg.id]
  db_subnet_group_name   = aws_db_subnet_group.rds_subnet_group.name


  tags = {
    Name = "Cloud-Platform-Kit-RDS-Database"
  }
}
resource "aws_ssm_parameter" "rds_db_connection_string" {
  name        = "/cloud-platform-kit/db/connection-string"
  description = "Full database connection string"
  type        = "SecureString"
  value       = "postgresql://${var.rds_database_username}:${var.rds_database_password}@${aws_db_instance.RDS_DB.address}:${aws_db_instance.RDS_DB.port}/${aws_db_instance.RDS_DB.db_name}"
  tags = {
    Name = "RDS Database full connection string"
  }
}