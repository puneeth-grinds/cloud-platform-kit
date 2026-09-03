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

resource "aws_eip" "elastic_ip" {
  domain = "vpc"


  tags = {
    Name = "Elastic-IP"
  }
}

resource "aws_nat_gateway" "nat_gateway" {
  allocation_id = aws_eip.elastic_ip.id
  subnet_id     = aws_subnet.public_1a.id
  depends_on    = [aws_internet_gateway.gateway]

  tags = {
    Name = "Nat-Gateway"
  }
}

resource "aws_route_table" "public-route-table" {
  vpc_id = aws_vpc.vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gateway.id
  }

  tags = {
    Name = "Public-Routing-Table"
  }

}

resource "aws_route_table" "private-route-table" {
  vpc_id = aws_vpc.vpc.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.nat_gateway.id
  }
  tags = {
    Name = "Private-Routing_Table"
  }
}

resource "aws_route_table_association" "public_route_table_association_1a" {
  subnet_id      = aws_subnet.public_1a.id
  route_table_id = aws_route_table.public-route-table.id
}

resource "aws_route_table_association" "public_route_table_association_1b" {
  subnet_id      = aws_subnet.public_1b.id
  route_table_id = aws_route_table.public-route-table.id
}

resource "aws_route_table_association" "private_route_table_association_1a" {
  subnet_id      = aws_subnet.private_1a.id
  route_table_id = aws_route_table.private-route-table.id
}

resource "aws_route_table_association" "private_route_table_association_1b" {
  subnet_id      = aws_subnet.private_1b.id
  route_table_id = aws_route_table.private-route-table.id
}