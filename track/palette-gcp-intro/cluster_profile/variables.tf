variable "sc_host" {
  description = "Spectro Cloud Endpoint"
  default     = "api.spectrocloud.com"
}

variable "sc_api_key" {
  description = "Spectro Cloud API key"
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


