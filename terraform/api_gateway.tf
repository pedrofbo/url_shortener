resource "aws_api_gateway_rest_api" "url_shortener" {
  name = "url_shortener"

  endpoint_configuration {
    types = [
      "REGIONAL",
    ]
  }
}