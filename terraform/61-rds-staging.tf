resource "aws_db_instance" "default" {
  identifier = "fieldkit-staging"

  tags {
    Name = "fieldkit-staging"
  }

  allocated_storage = 100
  engine            = "postgres"
  engine_version    = "9.6.1"
  instance_class    = "db.t2.small"

  name     = "fieldkit"
  username = "fieldkit"
  password = "kieldfit"

  db_subnet_group_name   = "${aws_db_subnet_group.fieldkit.name}"
  vpc_security_group_ids = ["${aws_security_group.postgresql-staging.id}"]
}
