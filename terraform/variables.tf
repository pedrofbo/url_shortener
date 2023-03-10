variable "profile" {
  type        = string
  description = "AWS Profile where the resources will be deployed."
}
variable "account_id" {
  type        = string
  description = "Account where the resources will be deployed."
}
variable "region" {
  type        = string
  description = "AWS Region where the resources will be deployed."
  default     = "us-east-1"
}
# TODO: build the executable at apply time.
variable "go_executable_path" {
  type        = string
  description = "Path to the (pre) built executable from the url shortener source code. Note that the executable must be named `url_shortener`."
  default     = "../url_shortener"
}
variable "base_endpoint" {
  type        = string
  description = "Base of the redirect URL returned by the API `create` endpoint. Can be either a custom domain or the invoke URL generated by API Gateway."
  default     = "https://short.pyoh.dev"
}
variable "default_redirect_endpoint" {
  type        = string
  description = "Default endpoint where the user will be redirected if a given shortened URL is not found in the database."
  default     = "https://fun.pyoh.dev"
}
variable "env_suffix" {
  type        = string
  description = "Suffix that will be appended at the end of resource names. Useful for multiple workspaces or having multiple instances of this module."
  default     = ""
}
