#Palette environment info
variable "sc_host" {
  description = "Spectro Cloud Endpoint"
  default     = "api.spectrocloud.com"
}

variable "sc_api_key" {
  description = "Palette api key"
  sensitive   = true
}

variable "sc_project_name" {
  description = "Spectro Cloud Project (e.g: Default)"
  default     = "Default"
}

#Cloud Account
#Variables used to create or reference a cloud account in Palette
variable "sc_cloud_account_name" {}
variable "aws_access_key_id" {}
variable "aws_secret_access_key" {
  sensitive    = true
}

variable "aws_ssh_key_name" {
  description  = "name of your key pair for the AWS region you are operating in"
}
variable "aws_key_pair" {
  description  = "your ssh public key"
}

#Cluster Profile
variable "sc_cp_profile_name" {}
variable "sc_cp_profile_description" {}
variable "sc_cp_profile_tags" {}
variable "sc_cp_cloud" {}
variable "sc_cp_type" {}

# Cluster
variable "sc_cluster_name" {}
variable "sc_cluster_tags" {}
variable "master_instance_type" {
  default = "t3a.large"  
}
variable "worker_instance_type" {
  default = "t3a.xlarge"  
}
variable "aws_region" {}
variable "master_aws_region_az" {}
variable "worker_aws_region_az" {}
variable "master-pool_node_count" {
  default = "1"
}
variable "worker-pool_node_count" {
  default = "1"
}



