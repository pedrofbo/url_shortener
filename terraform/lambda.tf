resource "aws_cloudwatch_log_group" "shorten" {
  # name = "/aws/lambda/${aws_lambda_function.ce_lambda.function_name}"
  name = "/aws/lambda/url_shortener__shorten"
}

resource "aws_cloudwatch_log_group" "redirect" {
  # name = "/aws/lambda/${aws_lambda_function.ce_lambda.function_name}"
  name = "/aws/lambda/url_shortener__redirect"
}
