resource "aws_lb" "lb" {
  name               = "cloudplatformkitlb"
  internal           = false
  load_balancer_type = "application"
  subnets            = [aws_subnet.public_1a.id, aws_subnet.private_1b.id]
  security_groups    = [aws_security_group.alb_sg.id]
}

resource "aws_alb_target_group" "api_alb_target" {
  name        = "api-gateway-target-group"
  vpc_id      = aws_vpc.vpc.id
  port        = 8080
  protocol    = "HTTP"
  target_type = "ip"

  health_check {
    path    = "/health"
    matcher = "200"
  }
}

resource "aws_lb_listener" "lb_listener" {
  port              = "8080"
  protocol          = "HTTP"
  load_balancer_arn = aws_lb.lb.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_alb_target_group.api_alb_target.arn
  }
}

resource "aws_ecs_cluster" "ecs_cluster" {
  name = "cloud-platform-kit"
}

resource "aws_service_discovery_private_dns_namespace" "dns_namespace" {
  vpc         = aws_vpc.vpc.id
  name        = "cloud-platform-kit.local"
  description = "Private DNS namespace for API gateway to find vul scanner"
}

resource "aws_service_discovery_service" "service_discovery_vul_scanner" {
  name = "vulnerability-scanner"
  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.dns_namespace.id

    dns_records {
      type = "A"
      ttl = 10
    }
    routing_policy = "MULTI-VALUE"
  }
}