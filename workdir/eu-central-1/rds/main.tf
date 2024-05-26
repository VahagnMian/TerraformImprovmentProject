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
  subnet_ids = getValueByKey("vpc", "public_subnet_ids")

  tags = {
    Name = "My DB subnet group"
  }
}
