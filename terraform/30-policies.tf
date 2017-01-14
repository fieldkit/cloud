data "aws_iam_policy_document" "fieldkit-server" {
  statement {
    actions = [
      "sqs:SendMessage",
      "sqs:DeleteMessage",
      "sqs:ReceiveMessage",
    ]
    resources = [
      "${aws_sqs_queue.fieldkit.arn}",
    ]
  }
}

resource "aws_iam_role" "fieldkit-server" {
  name = "fieldkit-server"
  assume_role_policy = <<EOF
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_instance_profile" "fieldkit-server" {
  name = "fieldkit-server"
  roles = ["${aws_iam_role.fieldkit-server.name}"]
}

resource "aws_iam_role_policy" "fieldkit-server" {
  name = "fieldkit-server"
  role = "${aws_iam_role.fieldkit-server.id}"
  policy = "${data.aws_iam_policy_document.fieldkit-server.json}"
}
