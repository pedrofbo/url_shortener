resource "aws_dynamodb_table" "entries" {
  name           = "url_shortener__entries${var.env_suffix}"
  hash_key       = "short_url"
  billing_mode   = "PROVISIONED"
  read_capacity  = 1
  write_capacity = 1

  attribute {
    name = "short_url"
    type = "S"
  }
}

resource "aws_dynamodb_table" "particles" {
  name           = "url_shortener__particles${var.env_suffix}"
  hash_key       = "particle"
  billing_mode   = "PROVISIONED"
  read_capacity  = 1
  write_capacity = 1

  attribute {
    name = "particle"
    type = "S"
  }
}
