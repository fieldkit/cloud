resource "aws_security_group" "ssh" {
  name = "ssh"
  description = "ssh"
  vpc_id = "${aws_vpc.fieldkit.id}"

  ingress {
      from_port = 22
      to_port = 22
      protocol = "tcp"
      cidr_blocks = ["24.103.10.214/32"]
  }

  egress {
      from_port = 0
      to_port = 0
      protocol = "-1"
      cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "fieldkit-server-elb" {
  name = "fieldkit-server-elb"
  description = "fieldkit-server-elb"
  vpc_id = "${aws_vpc.fieldkit.id}"

  ingress {
      from_port = 443
      to_port = 443
      protocol = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
      from_port = 0
      to_port = 0
      protocol = "-1"
      cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "fieldkit-server" {
  name = "fieldkit-server"
  description = "fieldkit-server"
  vpc_id = "${aws_vpc.fieldkit.id}"

  ingress {
      from_port = 80
      to_port = 80
      protocol = "tcp"
      security_groups = ["${aws_security_group.fieldkit-server-elb.id}"]
  }

  egress {
      from_port = 0
      to_port = 0
      protocol = "-1"
      cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "postgresql" {
  name = "postgresql"
  description = "postgresql"
  vpc_id = "${aws_vpc.fieldkit.id}"

  ingress {
    from_port = 5432
    to_port = 5432
    protocol = "tcp"
    security_groups = ["${aws_security_group.fieldkit-server.id}"]
  }
}

resource "aws_security_group" "redis" {
  name = "redis"
  description = "redis"
  vpc_id = "${aws_vpc.fieldkit.id}"

  ingress {
    from_port = 6379
    to_port = 6379
    protocol = "tcp"
    security_groups = ["${aws_security_group.fieldkit-server.id}"]
  }
}
