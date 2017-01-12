resource "aws_vpc" "fieldkit" {
  cidr_block = "10.0.0.0/16"
  tags {
    Name = "fieldkit"
  }
}

resource "aws_internet_gateway" "fieldkit" {
  vpc_id = "${aws_vpc.fieldkit.id}"
  tags {
    Name = "fieldkit"
  }
}

resource "aws_subnet" "fieldkit-a" {
	vpc_id = "${aws_vpc.fieldkit.id}"
	cidr_block = "10.0.0.0/18"
	availability_zone = "us-east-1a"
	map_public_ip_on_launch = true
	tags {
		Name = "fieldkit-a"
	}
}

resource "aws_subnet" "fieldkit-c" {
	vpc_id = "${aws_vpc.fieldkit.id}"
	cidr_block = "10.0.64.0/18"
	availability_zone = "us-east-1c"
	map_public_ip_on_launch = true
	tags {
		Name = "fieldkit-c"
	}
}

resource "aws_subnet" "fieldkit-d" {
	vpc_id = "${aws_vpc.fieldkit.id}"
	cidr_block = "10.0.128.0/18"
	availability_zone = "us-east-1d"
	map_public_ip_on_launch = true
	tags {
		Name = "fieldkit-d"
	}
}

resource "aws_subnet" "fieldkit-e" {
	vpc_id = "${aws_vpc.fieldkit.id}"
	cidr_block = "10.0.192.0/18"
	availability_zone = "us-east-1e"
	map_public_ip_on_launch = true
	tags {
		Name = "fieldkit-e"
	}
}

resource "aws_db_subnet_group" "fieldkit" {
  name = "fieldkit"
  description = "fieldkit"
  subnet_ids = ["${aws_subnet.fieldkit-a.id}", "${aws_subnet.fieldkit-c.id}", "${aws_subnet.fieldkit-d.id}", "${aws_subnet.fieldkit-e.id}"]
  tags {
    Name = "fieldkit"
  }
}

resource "aws_elasticache_subnet_group" "fieldkit" {
  name = "fieldkit"
  description = "fieldkit"
  subnet_ids = ["${aws_subnet.fieldkit-a.id}", "${aws_subnet.fieldkit-c.id}", "${aws_subnet.fieldkit-d.id}", "${aws_subnet.fieldkit-e.id}"]
}
