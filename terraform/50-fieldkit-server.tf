data "template_file" "fieldkit-server-a" {
  template = "${file("cloud-config/fieldkit-server.yaml")}"
  vars {
    hostname = "fieldkit-server-a"
  }
}

resource "aws_instance" "fieldkit-server-a" {
  ami = "ami-4d795c5a"
  instance_type = "t2.medium"
  subnet_id = "${aws_subnet.fieldkit-a.id}"
  associate_public_ip_address = true
  vpc_security_group_ids = ["${aws_security_group.ssh.id}", "${aws_security_group.fieldkit-server.id}"]
  user_data = "${data.template_file.fieldkit-server-a.rendered}"
  key_name = "fieldkit"
  iam_instance_profile = "${aws_iam_instance_profile.fieldkit-server.id}"
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
  name = "fieldkit-server-a.aws.fieldkit.org"
  type = "A"
  ttl = "60"
  records = ["${aws_instance.fieldkit-server-a.public_ip}"]
}

resource "aws_elb" "fieldkit" {
  name = "fieldkit"
  subnets = ["${aws_subnet.fieldkit-a.id}", "${aws_subnet.fieldkit-c.id}", "${aws_subnet.fieldkit-d.id}", "${aws_subnet.fieldkit-e.id}"]
  security_groups = ["${aws_security_group.fieldkit-server-elb.id}"]
  
  listener {
    instance_port = 80
    instance_protocol = "http"
    lb_port = 443
    lb_protocol = "https"
    ssl_certificate_id = "arn:aws:acm:us-east-1:582827299311:certificate/4a280325-b065-450b-b021-8862b52a7d3e"
  }

  health_check {
    healthy_threshold = 2
    unhealthy_threshold = 2
    timeout = 3
    target = "HTTP:80/"
    interval = 30
  }

  instances = ["${aws_instance.fieldkit-server-a.id}"]
  cross_zone_load_balancing = true

  tags {
    Name = "fieldkit"
  }
}
