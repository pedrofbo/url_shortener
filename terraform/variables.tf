variable "profile" {
  type        = string
  description = "AWS Profile where the resources will be deployed"
}

variable "region" {
  type        = string
  description = "AWS Region where the resources will be deployed"
  default     = "us-east-1"
}
