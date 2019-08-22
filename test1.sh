#!/bin/bash

dd if=/dev/urandom of=/tmp/100000 bs=100000 count=1

curl http://127.0.0.1:8080/ingestion --data-binary @/tmp/100000 \
     -H "Content-Type: application/vnd.fk.data+binary" \
     -H "Fk-DeviceId: ZWZVjVM1NTM0ICAg/xg/Og==" \
     -H "Fk-Blocks: 0,100" \
     -H "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImphY29iQGNvbnNlcnZpZnkub3JnIiwiZXhwIjoxNTY2NTc1ODgyLCJpYXQiOjE1NjY0ODk0ODIsInJlZnJlc2hfdG9rZW4iOiJzeEh0QnpaelF2bXlOVlNweTlEa0djM0wzb0EiLCJzY29wZXMiOiJhcGk6YWNjZXNzIiwic3ViIjoyfQ.__vatXV417Ecq5e9bfhao1WENy0yLkeNbCWj4IsOj5oALoUX-AuHWfL_XZqK-gX4_OA-ISeEKA4zhJscmU6Ddw"

rm -f /tmp/100000
