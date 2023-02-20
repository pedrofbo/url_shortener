resource "aws_dynamodb_table" "entries" {
  name         = "url_shortener__entries"
  hash_key     = "short_url"
  billing_mode = "PROVISIONED"

  attribute {
    name = "short_url"
    type = "S"
  }
}

resource "aws_dynamodb_table" "particles" {
  name         = "url_shortener__particles"
  hash_key     = "particle"
  billing_mode = "PROVISIONED"

  attribute {
    name = "particle"
    type = "S"
  }
}
