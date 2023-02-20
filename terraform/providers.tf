provider "aws" {
  profile = var.profile
  region  = var.region
  default_tags {
    tags = {
      "project"     = "url_shortener",
      "project_ref" = "https://github.com/pedrofbo/url_shortener",
    }
  }
}
