
resource "aws_security_group" "alb_sg" {
  name        = "allow on port 80"
  description = "Allow all traffic from internet on port 80"
  vpc_id      = aws_vpc.vpc.id
}

resource "aws_security_group" "ecs_sg" {
  name        = "allow alb"
  description = "Allow the ALB to reach the ECS application"
  vpc_id      = aws_vpc.vpc.id
}

resource "aws_security_group" "rds_sg" {
  name        = "allow rds"
  description = "Allow the rds from application"
  vpc_id      = aws_vpc.vpc.id
}
resource "aws_vpc_security_group_ingress_rule" "alb_sg_ingress_rule" {
  security_group_id = aws_security_group.alb_sg.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 80
  to_port           = 80
  ip_protocol       = "tcp"

}
resource "aws_vpc_security_group_egress_rule" "alb_sg_egress_rule" {
  security_group_id = aws_security_group.alb_sg.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 0
  to_port           = 0
  ip_protocol       = -1
}
resource "aws_vpc_security_group_ingress_rule" "ecs_sg_ingress_rule" {
  security_group_id            = aws_security_group.ecs_sg.id
  referenced_security_group_id = aws_security_group.alb_sg.id
  from_port                    = 8080
  to_port                      = 8080
  ip_protocol                  = "tcp"

}

resource "aws_vpc_security_group_ingress_rule" "ecs_self_sg_ingress_rule" {
  security_group_id            = aws_security_group.ecs_sg.id
  referenced_security_group_id = aws_security_group.ecs_sg.id
  from_port                    = 8081
  to_port                      = 8081
  ip_protocol                  = "tcp"
}

resource "aws_vpc_security_group_egress_rule" "ecs_sg_egress_rds" {
  security_group_id            = aws_security_group.ecs_sg.id
  referenced_security_group_id = aws_security_group.rds_sg.id
  from_port                    = 5432
  to_port                      = 5432
  ip_protocol                  = "tcp"
}