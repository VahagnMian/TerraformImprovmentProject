resource "aws_db_instance" "default" {
  allocated_storage    = ""
  db_subnet_group_name = aws_db_subnet_group.default.name
  db_name              = "mydb"
  engine               = "mysql"
  engine_version       = "5.7"
  instance_class       = "db.t3.micro"
  username             = "foo"
  password             = "foobarbaz"
  parameter_group_name = "default.mysql5.7"
  skip_final_snapshot  = true
}

resource "aws_db_subnet_group" "default" {
  name       = "main"
  subnet_ids = ["subnet-0ed3c309930d86048", "subnet-08efa53d05510ca49", "subnet-00db38c17caf7ab42"]

  tags = {
    Name = "My DB subnet group"
  }
}
