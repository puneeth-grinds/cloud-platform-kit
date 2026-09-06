resource "aws_lb" "lb" {
  name = "app_lb"
  internal = false
  load_balancer_type = "application"
}