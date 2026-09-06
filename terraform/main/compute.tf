resource "aws_lb" "lb" {
  name = "app_lb"
  internal = false
  load_balancer_type = "application"
  subnets = [ aws_subnet.public_1a, aws_subnet.private_1b ]
  security_groups = [ aws_security_group.alb_sg.id ]
}