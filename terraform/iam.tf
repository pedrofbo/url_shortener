resource "aws_iam_role" "lambda" {
  name = "UrlShortenerLambdaRole${var.env_suffix}"
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      },
    ]
  })
}

resource "aws_iam_role_policy" "url_shortener_lambda_permissions" {
  name   = "url_shortener-lambda-permissions${var.env_suffix}"
  role   = aws_iam_role.lambda.name
  policy = data.aws_iam_policy_document.url_shortener_lambda_permissions.json
}

data "aws_iam_policy_document" "url_shortener_lambda_permissions" {
  statement {
    effect = "Allow"
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      aws_cloudwatch_log_group.shorten.arn,
      "${aws_cloudwatch_log_group.shorten.arn}:*",
      aws_cloudwatch_log_group.redirect.arn,
      "${aws_cloudwatch_log_group.redirect.arn}:*"
    ]
  }
  statement {
    effect = "Allow"
    actions = [
      "dynamodb:BatchGetItem",
      "dynamodb:BatchWriteItem",
      "dynamodb:ConditionCheckItem",
      "dynamodb:PutItem",
      "dynamodb:DeleteItem",
      "dynamodb:GetItem",
      "dynamodb:Scan",
      "dynamodb:Query",
      "dynamodb:UpdateItem"
    ]
    resources = [
      aws_dynamodb_table.entries.arn,
      aws_dynamodb_table.particles.arn
    ]
  }
}
