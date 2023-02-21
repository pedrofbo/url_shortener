##########################
###  Lambda Functions  ###
##########################
data "archive_file" "lambda_zip" {
  type             = "zip"
  output_file_mode = "0644"
  source_file      = var.go_executable_path
  output_path      = "./url_shortener.zip"
}

resource "aws_lambda_function" "shorten" {
  function_name = "url_shortener__shorten${var.env_suffix}"
  filename      = data.archive_file.lambda_zip.output_path
  role          = aws_iam_role.lambda.arn
  handler       = "url_shortener"

  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
  runtime          = "go1.x"
  memory_size      = 512
  timeout          = 15

  environment {
    variables = {
      "API_HANDLER"          = "LAMBDA_SHORTEN"
      "BASE_ENDPOINT"        = var.base_endpoint
      "ENTRIES_TABLE_NAME"   = aws_dynamodb_table.entries.name
      "PARTICLES_TABLE_NAME" = aws_dynamodb_table.particles.name
    }
  }
}

resource "aws_lambda_function" "redirect" {
  function_name = "url_shortener__redirect${var.env_suffix}"
  filename      = data.archive_file.lambda_zip.output_path
  role          = aws_iam_role.lambda.arn
  handler       = "url_shortener"

  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
  runtime          = "go1.x"
  memory_size      = 512
  timeout          = 15

  environment {
    variables = {
      "API_HANDLER"               = "LAMBDA_REDIRECT"
      "DEFAULT_REDIRECT_ENDPOINT" = var.default_redirect_endpoint
      "ENTRIES_TABLE_NAME"        = aws_dynamodb_table.entries.name
      "PARTICLES_TABLE_NAME"      = aws_dynamodb_table.particles.name
    }
  }
}

#########################
###  Cloudwatch Logs  ###
#########################
resource "aws_cloudwatch_log_group" "shorten" {
  name = "/aws/lambda/${aws_lambda_function.shorten.function_name}"
}

resource "aws_cloudwatch_log_group" "redirect" {
  name = "/aws/lambda/${aws_lambda_function.redirect.function_name}"
}
