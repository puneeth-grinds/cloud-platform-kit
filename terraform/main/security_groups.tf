
resource "aws_security_group" "alb_sg" {
  name        = "allow on port 80"
  description = "Allow all traffic from internet on port 80"
  vpc_id      = aws_vpc.vpc.id
}