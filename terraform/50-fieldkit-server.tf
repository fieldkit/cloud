data "template_file" "fieldkit-server-a" {
  template = "${file("cloud-config/fieldkit-server.yaml")}"

  vars {
    hostname = "fieldkit-server-a"
  }
}

resource "aws_instance" "fieldkit-server-a" {
  ami                         = "ami-3b7f9e2d"
  instance_type               = "m4.large"
  subnet_id                   = "${aws_subnet.fieldkit-a.id}"
  associate_public_ip_address = true
  vpc_security_group_ids      = ["${aws_security_group.ssh.id}", "${aws_security_group.fieldkit-server.id}"]
  user_data                   = "${data.template_file.fieldkit-server-a.rendered}"
  key_name                    = "fieldkit"
  iam_instance_profile        = "${aws_iam_instance_profile.fieldkit-server.id}"

  root_block_device {
    volume_type = "gp2"
    volume_size = 100
  }

  tags {
    Name = "fieldkit-server-a"
  }
}

resource "aws_route53_record" "fieldkit-server-a" {
  zone_id = "Z116TCZ3RT5Z2K"
  name    = "fieldkit-server-a.aws.fieldkit.org"
  type    = "A"
  ttl     = "60"
  records = ["${aws_instance.fieldkit-server-a.public_ip}"]
}

resource "aws_alb" "fieldkit-server" {
  name            = "fieldkit-server"
  internal        = false
  security_groups = ["${aws_security_group.fieldkit-server-alb.id}"]
  subnets         = ["${aws_subnet.fieldkit-a.id}", "${aws_subnet.fieldkit-b.id}", "${aws_subnet.fieldkit-c.id}", "${aws_subnet.fieldkit-e.id}"]

  tags {
    Name = "fieldkit-server"
  }
}

resource "aws_alb_listener" "fieldkit-server" {
  load_balancer_arn = "${aws_alb.fieldkit-server.arn}"
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS-1-2-2017-01"
  certificate_arn   = "arn:aws:acm:us-east-1:582827299311:certificate/4a280325-b065-450b-b021-8862b52a7d3e"

  default_action {
    target_group_arn = "${aws_alb_target_group.fieldkit-server.arn}"
    type             = "forward"
  }
}

resource "aws_alb_target_group" "fieldkit-server" {
  name     = "fieldkit-server"
  port     = 80
  protocol = "HTTP"
  vpc_id   = "${aws_vpc.fieldkit.id}"

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    port                = 80
    path                = "/status"
    interval            = 5
  }
}

resource "aws_alb_target_group_attachment" "fieldkit-server-a" {
  target_group_arn = "${aws_alb_target_group.fieldkit-server.arn}"
  target_id        = "${aws_instance.fieldkit-server-a.id}"
  port             = 80
}

resource "aws_route53_record" "api-data" {
  zone_id = "Z116TCZ3RT5Z2K"
  name    = "api.data.fieldkit.org"
  type    = "A"

  alias {
    name                   = "${aws_alb.fieldkit-server.dns_name}"
    zone_id                = "${aws_alb.fieldkit-server.zone_id}"
    evaluate_target_health = false
  }
}
