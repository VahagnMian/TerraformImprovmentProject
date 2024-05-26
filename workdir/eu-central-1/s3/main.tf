resource "aws_s3_bucket" "bucket" {
  bucket = "miandevops-terraform-improvments"

  tags = {
    Name        = "My bucket"
    Environment = "Dev"
    Purpose = "Temporary"
  }
}