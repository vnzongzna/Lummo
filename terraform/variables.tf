variable "do_token" {
  type      = string
  sensitive = true
}

variable "region" {
  type    = string
  default = "tor1"
}

variable "k8s-version" {
  type    = string
  default = "1.21.10-do.0"
}
