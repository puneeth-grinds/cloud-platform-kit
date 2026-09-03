
resource "aws_security_group" "alb_sg" {
  name        = "allow on port 80"
  description = "Allow all traffic from internet on port 80"
  vpc_id      = aws_vpc.vpc.id
}

resource "aws_vpc_security_group_ingress_rule" "alb_sg_ingress_rule" {
  security_group_id = aws_security_group.alb_sg.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 80
  ip_protocol       = "tcp"
  to_port           = 80
}