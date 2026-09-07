resource "aws_lb" "lb" {
  name               = "app_lb"
  internal           = false
  load_balancer_type = "application"
  subnets            = [aws_subnet.public_1a, aws_subnet.private_1b]
  security_groups    = [aws_security_group.alb_sg.id]
}

resource "aws_alb_target_group" "api_alb_target" {
  name = "api-gateway-target-group"
  vpc_id = aws_vpc.vpc.id
  port = 8080
  protocol = "HTTP"
  target_type = "ip"

  health_check {
    path = "/health"
    matcher = "200"
  }
}