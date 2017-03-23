resource "aws_security_group" "fieldkit-server-staging-alb" {
  name        = "fieldkit-server-staging-alb"
  description = "fieldkit-server-staging-alb"
  vpc_id      = "${aws_vpc.fieldkit.id}"

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "fieldkit-server-staging" {
  name        = "fieldkit-server-staging"
  description = "fieldkit-server-staging"
  vpc_id      = "${aws_vpc.fieldkit.id}"

  ingress {
    from_port       = 80
    to_port         = 80
    protocol        = "tcp"
    security_groups = ["${aws_security_group.fieldkit-server-staging-alb.id}"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "postgresql-staging" {
  name        = "postgresql-staging"
  description = "postgresql-staging"
  vpc_id      = "${aws_vpc.fieldkit.id}"

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = ["${aws_security_group.fieldkit-server-staging.id}"]
  }
}
