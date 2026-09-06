
resource "aws_db_subnet_group" "rds_subnet_group" {
  name       = "RDS_Subnet"
  subnet_ids = [aws_subnet.private_1a, aws_subnet.private_1b]

  tags = {
    Name = "RDS subnet group"
  }
}

resource "aws_db_instance" "RDS_DB" {
  engine                 = "PostgreSQL"
  engine_version         = "16"
  instance_class         = "db.t3.micro"
  allocated_storage      = "20GB"
  multi_az               = false
  availability_zone      = var.availability_zone_1a
  db_name                = var.rds_database_name
  username               = var.rds_database_name
  password               = var.rds_database_password
  publicly_accessible    = false
  skip_final_snapshot    = true
  deletion_protection    = false
  vpc_security_group_ids = [aws_db_subnet_group.rds_subnet_group.id]

  tags = {
    Name = "Cloud-Platform-Kit-RDS-Database"
  }
}

resource "aws_ssm_parameter" "db_endpoint" {
  name        = "/cloud-platform-kit/db/connection-string"
  description = "Database Endpoint"
  type        = "secureString"
  value       = aws_db_instance.RDS_DB.endpoint
  tags = {
    Name = "RDS Database endpoint ssm parameter"
  }
}

resource "aws_ssm_parameter" "db_port" {
  name        = "/cloud-platform-kit/db/connection-string"
  description = "Database Port"
  type        = "secureString"
  value       = aws_db_instance.RDS_DB.port
  tags = {
    Name = "RDS Database port ssm parameter"
  }
}

resource "aws_ssm_parameter" "db_name" {
  name        = "/cloud-platform-kit/db/connection-string"
  description = "Database Port"
  type        = "secureString"
  value       = aws_db_instance.RDS_DB.db_name
  tags = {
    Name = "RDS Database name ssm parameter"
  }
}

resource "aws_ssm_parameter" "db_username" {
  name        = "/cloud-platform-kit/db/connection-string"
  description = "Database Port"
  type        = "secureString"
  value       = aws_db_instance.RDS_DB.username
  tags = {
    Name = "RDS Database name ssm username"
  }
}

resource "aws_ssm_parameter" "db_password" {
  name        = "/cloud-platform-kit/db/connection-string"
  description = "Database Port"
  type        = "secureString"
  value       = aws_db_instance.RDS_DB.password
  tags = {
    Name = "RDS Database name ssm password"
  }
}