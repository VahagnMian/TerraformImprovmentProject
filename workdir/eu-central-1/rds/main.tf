resource "aws_db_instance" "default" {
  allocated_storage    = "150"
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
  subnet_ids = ["subnet-06756858d8b44da6a", "subnet-028866daba4f44c24", "subnet-0eff510bd9a782b31"]

  tags = {
    Name = "My DB subnet group"
  }
}
