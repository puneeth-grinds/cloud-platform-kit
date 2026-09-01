resource "aws_vpc" "vpc" {
  cidr_block           = var.vpc_cidr_block
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name = "cloud-platform-kit-vpc"
  }
}
resource "aws_internet_gateway" "gateway" {
  vpc_id = aws_vpc.vpc.id
  tags = {
    Name = "cloud-platform-kit-internet-gateway"
  }
}

resource "aws_subnet" "public_1a" {
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = var.public_1a
  availability_zone = var.availability_zone_1a

  tags = {
    Name = "public-subnet-1a"
  }
}

resource "aws_subnet" "public_1b" {
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = var.public_1b
  availability_zone = var.availability_zone_1b

  tags = {
    Name = "public-subnet-1b"
  }
}

resource "aws_subnet" "private_1a" {
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = var.private_1a
  availability_zone = var.availability_zone_1a
  tags = {
    Name = "private-subnet-1a"
  }
}

resource "aws_subnet" "private_1b" {
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = var.private_1b
  availability_zone = var.availability_zone_1b
  tags = {
    Name = "private-subnet-1b"
  }
}