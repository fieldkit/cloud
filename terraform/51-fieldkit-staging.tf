data "template_file" "fieldkit-server-staging-a" {
  template = "${file("cloud-config/fieldkit-server.yaml")}"

  vars {
    hostname = "fieldkit-server-staging-a"
  }
}

resource "aws_instance" "fieldkit-server-staging-a" {
  ami                         = "ami-3b7f9e2d"
  instance_type               = "t2.medium"
  subnet_id                   = "${aws_subnet.fieldkit-a.id}"
  associate_public_ip_address = true
  vpc_security_group_ids      = ["${aws_security_group.ssh.id}", "${aws_security_group.fieldkit-server-staging.id}"]
  user_data                   = "${data.template_file.fieldkit-server-staging-a.rendered}"
  key_name                    = "fieldkit"
  iam_instance_profile        = "${aws_iam_instance_profile.fieldkit-server-staging.id}"

  root_block_device {
    volume_type = "gp2"
    volume_size = 100
  }

  tags {
    Name = "fieldkit-server-staging-a"
  }
}

resource "aws_alb" "fieldkit-server-staging" {
  name            = "fieldkit-server-staging"
  internal        = false
  security_groups = ["${aws_security_group.fieldkit-server-staging-alb.id}"]
  subnets         = ["${aws_subnet.fieldkit-a.id}", "${aws_subnet.fieldkit-b.id}", "${aws_subnet.fieldkit-c.id}", "${aws_subnet.fieldkit-e.id}"]

  tags {
    Name = "fieldkit-server-staging"
  }
}

resource "aws_alb_listener" "fieldkit-server-staging" {
  load_balancer_arn = "${aws_alb.fieldkit-server-staging.arn}"
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS-1-2-2017-01"
  certificate_arn   = "arn:aws:acm:us-east-1:582827299311:certificate/30a77491-8d27-4ad5-af2d-b134ed9c73f9"

  default_action {
    target_group_arn = "${aws_alb_target_group.fieldkit-server-staging.arn}"
    type             = "forward"
  }
}

resource "aws_alb_target_group" "fieldkit-server-staging" {
  name     = "fieldkit-server-staging"
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

resource "aws_alb_target_group_attachment" "fieldkit-server-staging-a" {
  target_group_arn = "${aws_alb_target_group.fieldkit-server-staging.arn}"
  target_id        = "${aws_instance.fieldkit-server-staging-a.id}"
  port             = 80
}

resource "aws_route53_record" "fieldkit-server-staging-a" {
  zone_id = "Z1P1FTADK9BR86"
  name    = "fieldkit-server-staging-a.aws.fieldkit.team"
  type    = "A"
  ttl     = "60"
  records = ["${aws_instance.fieldkit-server-staging-a.public_ip}"]
}

resource "aws_route53_record" "frontend-staging" {
  zone_id = "Z1P1FTADK9BR86"
  name    = "fieldkit.team"
  type    = "A"

  alias {
    name                   = "${aws_alb.fieldkit-server-staging.dns_name}"
    zone_id                = "${aws_alb.fieldkit-server-staging.zone_id}"
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "api-data-staging" {
  zone_id = "Z1P1FTADK9BR86"
  name    = "api.fieldkit.team"
  type    = "A"

  alias {
    name                   = "${aws_alb.fieldkit-server-staging.dns_name}"
    zone_id                = "${aws_alb.fieldkit-server-staging.zone_id}"
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "wildcard-staging" {
  zone_id = "Z1P1FTADK9BR86"
  name    = "*.fieldkit.team"
  type    = "A"

  alias {
    name                   = "${aws_alb.fieldkit-server-staging.dns_name}"
    zone_id                = "${aws_alb.fieldkit-server-staging.zone_id}"
    evaluate_target_health = false
  }
}
