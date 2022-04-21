variable "sc_host" {
  description = "Spectro Cloud Endpoint"
  default     = "api.spectrocloud.com"
}

variable "sc_username" {
  description = "Spectro Cloud Username"
}

variable "sc_password" {
  description = "Spectro Cloud Password"
  sensitive   = true
}

variable "sc_api_key" {
  description = "Palette API key"
  sensitive   = true
}

variable "sc_project_name" {
  description = "Spectro Cloud Project (e.g: Default)"
  default     = "Default"
}

#Cluster Profile
variable "cluster_profile_name" {}
variable "cluster_profile_description" {}
variable "cluster_profile_tags" {}
variable "cloud" {}
variable "type" {}


