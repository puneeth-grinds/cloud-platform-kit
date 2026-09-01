resource "aws_vpc" "vpc" {
  cidr_block           = var.vpc_cidr_block
  enable_dns_support   = true
  enable_dns_hostnames = true
}

resource "aws_internet_gateway" "gateway" {
  vpc_id = aws_vpc.vpc.id
}

resource "aws_subnet" "public_1a" {
  vpc_id     = aws_vpc.vpc.id
  cidr_block = var.public_1a
}

resource "aws_subnet" "public_1b" {
  vpc_id     = aws_vpc.vpc.id
  cidr_block = var.public_1b
}

resource "aws_subnet" "private_1a" {
  vpc_id     = aws_vpc.vpc.id
  cidr_block = var.private_1a
}

resource "aws_subnet" "private_1b" {
  vpc_id     = aws_vpc.vpc.id
  cidr_block = var.private_1b
}