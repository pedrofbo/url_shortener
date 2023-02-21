resource "aws_api_gateway_rest_api" "url_shortener" {
  name = "url_shortener"

  endpoint_configuration {
    types = [
      "REGIONAL",
    ]
  }
}

##############################
###  Shorten URL endpoint  ###
##############################

resource "aws_api_gateway_resource" "create" {
  path_part   = "create"
  rest_api_id = aws_api_gateway_rest_api.url_shortener.id
  parent_id   = aws_api_gateway_rest_api.url_shortener.root_resource_id
}

resource "aws_api_gateway_method" "create__post" {
  rest_api_id   = aws_api_gateway_rest_api.url_shortener.id
  resource_id   = aws_api_gateway_resource.create.id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "create__post__lambda" {
  rest_api_id             = aws_api_gateway_rest_api.url_shortener.id
  resource_id             = aws_api_gateway_resource.create.id
  http_method             = "POST"
  type                    = "AWS_PROXY"
  integration_http_method = "POST"
  content_handling        = "CONVERT_TO_TEXT"
  uri                     = aws_lambda_function.shorten.invoke_arn
}

resource "aws_lambda_permission" "api_gateway_invoke_lambda" {
  function_name = aws_lambda_function.shorten.arn
  action        = "lambda:InvokeFunction"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "arn:aws:execute-api:${var.region}:${local.account_id}:${aws_api_gateway_rest_api.url_shortener.id}/*/${aws_api_gateway_method.create__post.http_method}${aws_api_gateway_resource.create.path}"
}
