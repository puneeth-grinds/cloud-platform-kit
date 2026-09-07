resource "aws_lb" "lb" {
  name               = "cloudplatformkitlb"
  internal           = false
  load_balancer_type = "application"
  subnets            = [aws_subnet.public_1a.id, aws_subnet.public_1b.id]
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
  port              = "80"
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
      ttl  = 10
    }
    routing_policy = "MULTIVALUE"
  }
}

resource "aws_cloudwatch_log_group" "cw_log_group_ecs" {
  name              = "/ecs/api-gateway"
  retention_in_days = 7
}

resource "aws_cloudwatch_log_group" "cw_log_group_vul" {
  name              = "/ecs/vulnerability-scanner"
  retention_in_days = 7
}

resource "aws_ecs_task_definition" "ecs_task_apigateway" {
  family                   = "api-gateway"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn
  task_role_arn            = aws_iam_role.api_gateway_task.arn

  container_definitions = jsonencode([
    {
      name      = "api-gateway"
      image     = "public.ecr.aws/docker/library/busybox:1.38.0"
      essential = true
      command = [
        "sh",
        "-c",
        "mkdir -p /www && echo ok > /www/health && exec httpd -f -p 8080 -h /www"
      ]
      portMappings = [
        {
          containerPort = 8080
          protocol      = "tcp"
        }
      ]
      secrets = [
        {
          name      = "DATABASE_URL"
          valueFrom = aws_ssm_parameter.rds_db_connection_string.arn
        }
      ]
      environment = [
        {
          name  = "SCANNER_URL"
          value = "http://vulnerability-scanner.cloud-platform-kit.local:8081"
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.cw_log_group_ecs.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])
}

resource "aws_ecs_task_definition" "ecs_task_vulscanner" {
  family                   = "vulnerability-scanner"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn
  task_role_arn            = aws_iam_role.vulnerability_scanner_task.arn

  container_definitions = jsonencode([
    {
      name      = "vulnerability-scanner"
      image     = "public.ecr.aws/docker/library/busybox:1.38.0"
      essential = true
      command = [
        "sh",
        "-c",
        "mkdir -p /www && echo ok > /www/health && exec httpd -f -p 8081 -h /www"
      ]
      portMappings = [
        {
          containerPort = 8081
          protocol      = "tcp"
        }
      ]
      secrets = [
        {
          name      = "DATABASE_URL"
          valueFrom = aws_ssm_parameter.rds_db_connection_string.arn
        }
      ]
      environment = [
        {
          name  = "S3_BUCKET"
          value = aws_s3_bucket.s3_storage.bucket
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.cw_log_group_vul.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])

}

resource "aws_ecs_service" "ecs_service_apigateway" {
  name            = "api-gateway"
  cluster         = aws_ecs_cluster.ecs_cluster.id
  task_definition = aws_ecs_task_definition.ecs_task_apigateway.arn
  desired_count   = 1
  launch_type     = "FARGATE"
  depends_on      = [aws_lb_listener.lb_listener]
  network_configuration {
    subnets          = [aws_subnet.private_1a.id, aws_subnet.private_1b.id]
    security_groups  = [aws_security_group.ecs_sg.id]
    assign_public_ip = false
  }
  load_balancer {
    target_group_arn = aws_alb_target_group.api_alb_target.arn
    container_name   = "api-gateway"
    container_port   = 8080
  }
  health_check_grace_period_seconds = 60
}

resource "aws_ecs_service" "ecs_service_vulscanner" {
  name            = "vulnerability-scanner"
  cluster         = aws_ecs_cluster.ecs_cluster.id
  task_definition = aws_ecs_task_definition.ecs_task_vulscanner.arn
  desired_count   = 1
  launch_type     = "FARGATE"
  depends_on      = [aws_lb_listener.lb_listener]
  network_configuration {
    subnets          = [aws_subnet.private_1a.id, aws_subnet.private_1b.id]
    security_groups  = [aws_security_group.ecs_sg.id]
    assign_public_ip = false
  }
  service_registries {
    registry_arn = aws_service_discovery_service.service_discovery_vul_scanner.arn
  }
}